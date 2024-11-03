package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
)

type TestContext struct {
	v map[string]interface{}
}

const sessionHeader = "x-session"

func (t *TestContext) Bind(v interface{}) error {
	*v.(*model.Todo) = model.Todo{Title: "sleep"}
	return nil
}
func (t *TestContext) JSON(code int, v interface{}) { t.v = v.(map[string]interface{}) }
func (t *TestContext) Log(string) logger.ILogDetail    { return logger.New(slog.Default(), "", nil) }
func (t *TestContext) Get(string) interface{}       { return nil }
func (t *TestContext) TransactionID() string        { return "" }
func (t *TestContext) Param(string) string          { return "" }
func (t *TestContext) Query(string) string          { return "" }
func (t *TestContext) Incoming() map[string]interface{} {
	return map[string]interface{}{}
}

type TestDB struct{}

func (*TestDB) Create(*model.Todo) error    { return nil }
func (*TestDB) List() ([]model.Todo, error) { return nil, nil }
func (*TestDB) Delete(id string) error      { return nil }

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set(sessionHeader, "test-ping")
	w := httptest.NewRecorder()
	c := &TestContext{}

	PingHandler(c)

	want := "pong"

	if c.v["message"] != want {
		t.Errorf("want %s, got %s", want, c.v["message"])
	}

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got code %d", http.StatusOK, w.Code)
	}

}

func TestX(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/x", nil)
	req.Header.Set(sessionHeader, "test-x")
	w := httptest.NewRecorder()
	c := &TestContext{}

	X(c)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set(sessionHeader, "test-health")
	w := httptest.NewRecorder()
	c := &TestContext{}

	Healthz(c)

	want := "ok"

	if c.v["status"] != want {
		t.Errorf("want %s, got %s", want, c.v["status"])
	}

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}
}

func TestTransfer(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/transfer/1", nil)
	req.Header.Set(sessionHeader, "test-")
	w := httptest.NewRecorder()
	c := &TestContext{}

	Transfer(c)

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}
}
