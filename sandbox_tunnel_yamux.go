package langsmith

import (
	"encoding/binary"
	"errors"
	"io"
	"sync"
	"time"
)

type yamuxSession struct {
	conn      io.ReadWriteCloser
	streams   map[uint32]*yamuxStream
	nextID    uint32
	mu        sync.Mutex
	writeMu   sync.Mutex
	closed    bool
	closeOnce sync.Once
}

func newYamuxSession(conn io.ReadWriteCloser) *yamuxSession {
	s := &yamuxSession{
		conn:    conn,
		streams: make(map[uint32]*yamuxStream),
		nextID:  1,
	}
	go s.readLoop()
	go s.keepAliveLoop()
	return s
}

func (s *yamuxSession) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *yamuxSession) openStream() (*yamuxStream, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, errors.New("yamux session is closed")
	}
	id := s.nextID
	s.nextID += 2
	stream := newYamuxStream(id, s)
	s.streams[id] = stream
	s.mu.Unlock()

	if err := s.sendFrame(yamuxTypeWindowUpdate, yamuxFlagSYN, id, 0); err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *yamuxSession) close() {
	s.closeOnce.Do(func() {
		s.mu.Lock()
		s.closed = true
		for _, stream := range s.streams {
			stream.receiveRST()
		}
		s.mu.Unlock()
		_ = s.sendFrame(yamuxTypeGoAway, 0, 0, 0)
		_ = s.conn.Close()
	})
}

func (s *yamuxSession) readLoop() {
	header := make([]byte, yamuxHeaderSize)
	for {
		if _, err := io.ReadFull(s.conn, header); err != nil {
			s.close()
			return
		}
		msgType := header[1]
		flags := binary.BigEndian.Uint16(header[2:4])
		streamID := binary.BigEndian.Uint32(header[4:8])
		length := binary.BigEndian.Uint32(header[8:12])

		switch msgType {
		case yamuxTypeData:
			payload := make([]byte, length)
			if length > 0 {
				if _, err := io.ReadFull(s.conn, payload); err != nil {
					s.close()
					return
				}
			}
			s.handleData(flags, streamID, payload)
		case yamuxTypeWindowUpdate:
			s.handleWindowUpdate(flags, streamID, length)
		case yamuxTypePing:
			if flags&yamuxFlagSYN != 0 {
				_ = s.sendFrame(yamuxTypePing, yamuxFlagACK, 0, length)
			}
		case yamuxTypeGoAway:
			s.close()
			return
		}
	}
}

func (s *yamuxSession) handleData(flags uint16, streamID uint32, payload []byte) {
	stream := s.getStream(streamID)
	if stream == nil {
		return
	}
	if len(payload) > 0 {
		stream.receiveData(payload)
	}
	if flags&yamuxFlagFIN != 0 {
		stream.receiveFIN()
	}
	if flags&yamuxFlagRST != 0 {
		stream.receiveRST()
	}
}

func (s *yamuxSession) handleWindowUpdate(flags uint16, streamID uint32, delta uint32) {
	stream := s.getStream(streamID)
	if stream == nil {
		return
	}
	if delta > 0 {
		stream.updateSendWindow(delta)
	}
	if flags&yamuxFlagFIN != 0 {
		stream.receiveFIN()
	}
	if flags&yamuxFlagRST != 0 {
		stream.receiveRST()
	}
}

func (s *yamuxSession) keepAliveLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	var ping uint32
	for range ticker.C {
		if s.isClosed() {
			return
		}
		ping++
		if err := s.sendFrame(yamuxTypePing, yamuxFlagSYN, 0, ping); err != nil {
			s.close()
			return
		}
	}
}

func (s *yamuxSession) getStream(streamID uint32) *yamuxStream {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.streams[streamID]
}

func (s *yamuxSession) sendFrame(msgType byte, flags uint16, streamID uint32, length uint32) error {
	var header [yamuxHeaderSize]byte
	header[0] = yamuxVersion
	header[1] = msgType
	binary.BigEndian.PutUint16(header[2:4], flags)
	binary.BigEndian.PutUint32(header[4:8], streamID)
	binary.BigEndian.PutUint32(header[8:12], length)
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_, err := s.conn.Write(header[:])
	return err
}

func (s *yamuxSession) sendData(streamID uint32, data []byte) error {
	var header [yamuxHeaderSize]byte
	header[0] = yamuxVersion
	header[1] = yamuxTypeData
	binary.BigEndian.PutUint32(header[4:8], streamID)
	binary.BigEndian.PutUint32(header[8:12], uint32(len(data)))
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if _, err := s.conn.Write(header[:]); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}

func (s *yamuxSession) sendWindowUpdate(streamID uint32, delta uint32) error {
	return s.sendFrame(yamuxTypeWindowUpdate, 0, streamID, delta)
}

type yamuxStream struct {
	id      uint32
	session *yamuxSession

	recvMu     sync.Mutex
	recvCond   *sync.Cond
	recvBuf    []byte
	recvClosed bool
	recvErr    bool
	recvWindow uint32

	sendMu     sync.Mutex
	sendCond   *sync.Cond
	sendClosed bool
	sendWindow uint32
}

func newYamuxStream(id uint32, session *yamuxSession) *yamuxStream {
	stream := &yamuxStream{
		id:         id,
		session:    session,
		recvWindow: yamuxInitialWindowSize,
		sendWindow: yamuxInitialWindowSize,
	}
	stream.recvCond = sync.NewCond(&stream.recvMu)
	stream.sendCond = sync.NewCond(&stream.sendMu)
	return stream
}

func (s *yamuxStream) Read(p []byte) (int, error) {
	s.recvMu.Lock()
	for len(s.recvBuf) == 0 && !s.recvClosed && !s.recvErr {
		s.recvCond.Wait()
	}
	if s.recvErr && len(s.recvBuf) == 0 {
		s.recvMu.Unlock()
		return 0, io.ErrUnexpectedEOF
	}
	if len(s.recvBuf) == 0 {
		s.recvMu.Unlock()
		return 0, io.EOF
	}
	n := copy(p, s.recvBuf)
	s.recvBuf = s.recvBuf[n:]
	consumed := yamuxInitialWindowSize - s.recvWindow
	if consumed >= yamuxInitialWindowSize/2 {
		s.recvWindow += consumed
		s.recvMu.Unlock()
		_ = s.session.sendWindowUpdate(s.id, consumed)
		return n, nil
	}
	s.recvMu.Unlock()
	return n, nil
}

func (s *yamuxStream) Write(p []byte) (int, error) {
	return s.write(p)
}

func (s *yamuxStream) write(p []byte) (int, error) {
	offset := 0
	for offset < len(p) {
		s.sendMu.Lock()
		for s.sendWindow == 0 && !s.sendClosed {
			s.sendCond.Wait()
		}
		if s.sendClosed {
			s.sendMu.Unlock()
			return offset, io.ErrClosedPipe
		}
		chunk := min(len(p)-offset, int(s.sendWindow))
		s.sendWindow -= uint32(chunk)
		s.sendMu.Unlock()

		if err := s.session.sendData(s.id, p[offset:offset+chunk]); err != nil {
			return offset, err
		}
		offset += chunk
	}
	return len(p), nil
}

func (s *yamuxStream) Close() error {
	s.close()
	return nil
}

func (s *yamuxStream) close() {
	s.sendMu.Lock()
	if !s.sendClosed {
		s.sendClosed = true
		_ = s.session.sendFrame(yamuxTypeData, yamuxFlagFIN, s.id, 0)
	}
	s.sendMu.Unlock()
	s.recvMu.Lock()
	s.recvClosed = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveData(data []byte) {
	s.recvMu.Lock()
	s.recvBuf = append(s.recvBuf, data...)
	s.recvWindow -= uint32(len(data))
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveFIN() {
	s.recvMu.Lock()
	s.recvClosed = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveRST() {
	s.recvMu.Lock()
	s.recvErr = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
	s.sendMu.Lock()
	s.sendClosed = true
	s.sendCond.Broadcast()
	s.sendMu.Unlock()
}

func (s *yamuxStream) updateSendWindow(delta uint32) {
	s.sendMu.Lock()
	s.sendWindow += delta
	s.sendCond.Broadcast()
	s.sendMu.Unlock()
}
