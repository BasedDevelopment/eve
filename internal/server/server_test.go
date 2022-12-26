package server

import (
	"net/http"
	"testing"
)

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get("http://localhost:3000/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected OK; got %v", resp.Status)
	}
}
