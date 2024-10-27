package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/todo"
)

type TestContext struct {
	v map[string]interface{}
}

func (t *TestContext) Bind(v interface{}) error {
	*v.(*todo.Todo) = todo.Todo{Title: "sleep"}
	return nil
}
func (t *TestContext) JSON(code int, v interface{}) {
	t.v = v.(map[string]interface{})
}
func (t *TestContext) Log() *logger.Logger {
	return logger.New(slog.Default(), nil)
}
func (t *TestContext) Get(string) interface{} { return nil }
func (t *TestContext) TransactionID() string  { return "" }
func (t *TestContext) Param(string) string    { return "" }

type TestDB struct{}

func (*TestDB) Create(*todo.Todo) error { return nil }

func (*TestDB) List() ([]todo.Todo, error) { return nil, nil }

func (*TestDB) Delete(id int) error { return nil }

func TestPing(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "http://localhost:8080/ping", nil)
	req.Header.Set("x-session", "test")
	w := httptest.NewRecorder()
	c := &TestContext{}
	// c.Request = req

	PingHandler(c)

	want := "pong"

	if c.v["message"] != want {
		t.Errorf("want %s, got %s", want, c.v["message"])
	}

	if w.Code != http.StatusOK {
		t.Errorf("want %d, got %d", http.StatusOK, w.Code)
	}

}
