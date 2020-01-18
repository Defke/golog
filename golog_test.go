package golog

import "testing"

func TestA(t *testing.T) {
	LoadConfig()
	Info("aaa", "bbbb")
}
