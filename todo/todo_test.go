package todo

import (
	"log/slog"
	"testing"

	"github.com/sing3demons/todoapi/logger"
)

type TestContext struct {
	v map[string]interface{}
}

func (t *TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{Title: "sleep"}
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

func (*TestDB) Create(*Todo) error { return nil }

func (*TestDB) List() ([]Todo, error) { return nil, nil }

func (*TestDB) Delete(id int) error { return nil }

func TestCreateTodoNotAllowSleep(t *testing.T) {
	handler := NewTodoHandler(&TestDB{})
	c := &TestContext{}
 	handler.NewTask(c)

	want := "not allowed"
	if c.v["error"] != want {
		t.Errorf("want %s got %s", want, c.v["error"])
	}

}
