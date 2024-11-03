package todo

import (
	"log/slog"
	"testing"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
	"github.com/sing3demons/todoapi/store"
)

type TestContext struct {
	v map[string]interface{}
}

func (t *TestContext) Bind(v interface{}) error {
	*v.(*model.Todo) = model.Todo{Title: "sleep"}
	return nil
}
func (t *TestContext) JSON(code int, v interface{}) {
	t.v = v.(map[string]interface{})
}
func (t *TestContext) Log(string) logger.ILogDetail {
	return logger.New(slog.Default(), "", nil)
}
func (t *TestContext) Get(string) interface{} { return nil }
func (t *TestContext) TransactionID() string  { return "" }
func (t *TestContext) Param(string) string    { return "" }
func (t *TestContext) Query(string) string    { return "" }

type TestDB struct{}

func (*TestDB) Create(*model.Todo) error { return nil }

func (*TestDB) List(store.FindOption) ([]model.Todo, error) { return nil, nil }

func (*TestDB) Delete(id string) error { return nil }
func (*TestDB) FindOne(id string) (*model.Todo, error) {
	return &model.Todo{Title: "sleep"}, nil
}

// func TestCreateTodo(t *testing.T) {
// 	handler := NewTodoHandler(&TestDB{})
// 	c := &TestContext{}
// 	handler.NewTask(c)

// 	want := "new task"
// 	if c.v["node"] != want {
// 		t.Errorf("want %s got %s", want, c.v["node"])
// 	}

// }

func TestCreateTodoNotAllowSleep(t *testing.T) {
	handler := NewTodoHandler(&TestDB{})
	c := &TestContext{}
	handler.NewTask(c)

	want := "not allowed"
	if c.v["error"] != want {
		t.Errorf("want %s got %s", want, c.v["error"])
	}

}
