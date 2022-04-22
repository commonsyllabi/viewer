package db

import (
	"testing"
	"time"
)

func TestConnect(t *testing.T) {

	time.Sleep(1 * time.Second)
	err := Connect("postgres", "postgres", "postgres", "localhost")

	if err != nil {
		t.Error(err)
	}
}
