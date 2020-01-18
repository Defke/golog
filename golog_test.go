package golog

import "testing"

func TestA(t *testing.T) {
	log, _ := LoadConfig()
	log.Info("aaa", "bbbb")
}
