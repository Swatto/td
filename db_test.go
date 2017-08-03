package main

import (
	"fmt"
	"os"
	"path"
	"testing"
	"time"
)

func TestWhenFileDoesntExist(t *testing.T) {
	cwd, _ := os.Getwd()
	extra := fmt.Sprint("/TODOtestingFOLDER/", time.Now().Format("20060102150405"))
	os.Setenv(ENVDBPATH, path.Join(cwd, extra))
	db, _ := NewDataStore()
	if db.Check() == nil {
		t.Errorf("Expected database check to return error, but it didn't.")
	}
	os.Unsetenv(ENVDBPATH)
}
