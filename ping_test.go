package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestPing(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	req.Header.Set("x-session", "test")
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	PingHandler(c)

	want := `{"message":"pong"}`

	if w.Body.String() != want {
		t.Errorf("want %s, got %s", want, w.Body.String())
	}

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}

}
