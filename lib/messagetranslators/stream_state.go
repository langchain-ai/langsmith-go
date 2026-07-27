package messagetranslators

func cloneMap[K comparable, V any](source map[K]V) map[K]V {
	cloned := make(map[K]V, len(source))
	for key, value := range source {
		cloned[key] = value
	}
	return cloned
}

func cloneBlockMap(source map[int]*anthropicBlock) map[int]*anthropicBlock {
	cloned := make(map[int]*anthropicBlock, len(source))
	for index, block := range source {
		copyBlock := *block
		cloned[index] = &copyBlock
	}
	return cloned
}
