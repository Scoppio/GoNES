package main

type Memory64k struct {
	mem [64 * 1024]byte
}

func (m *Memory64k) Reset() {
	m.mem = [64 * 1024]byte{}
}
