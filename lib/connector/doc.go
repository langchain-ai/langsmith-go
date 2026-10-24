// Package connector runs a LangSmith private network connector: the process
// that serves a tunnel or a relay from inside a network LangSmith cannot reach
// directly. It dials LangSmith over outbound WebSockets and carries the
// connections LangSmith opens through them, so the network needs no inbound
// access.
//
//	config, err := connector.NewConfig(connector.Settings{
//		ConnectionID: connectionID,
//		APIKey:       os.Getenv("LANGSMITH_API_KEY"),
//		Targets:      map[string]string{"app": "http://localhost:8080"},
//	})
//	if err != nil {
//		return err
//	}
//	client, err := connector.NewClient(config)
//	if err != nil {
//		return err
//	}
//	return client.Run(ctx)
//
// A tunnel serves named targets. A relay (Settings.Relay) instead dials the
// destinations LangSmith names, within its allow and deny rules. Every flow
// refuses cloud metadata endpoints.
package connector
