package main

import "testing"

// TODO: Need to actually test not just for errors
func TestGetDiffFromCmd(t *testing.T) {
	_, err := GetDiffFromCmd()
	if err != nil {
		t.Error(err)
	}
}
