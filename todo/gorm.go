package todo

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (g *GormStore) Create(todo *Todo) error {
	return g.db.Create(todo).Error
}

func (g *GormStore) List() ([]Todo, error) {
	var todos []Todo
	r := g.db.Find(&todos)
	if err := r.Error; err != nil {
		return nil, err
	}
	return todos, nil
}

func (g *GormStore) Get(id int) (*Todo, error) {
	var todo Todo
	r := g.db.First(&todo, id)
	if err := r.Error; err != nil {
		return nil, err
	}
	return &todo, nil
}

func (g *GormStore) Update(todo *Todo) error {
	return g.db.Save(todo).Error
}

func (g *GormStore) Delete(id int) error {
	return g.db.Delete(&Todo{}, id).Error
}

func (g *GormStore) Find(todos []Todo) ([]Todo, error) {
	r := g.db.Find(&todos)
	if err := r.Error; err != nil {
		return nil, err
	}
	return todos, nil
}