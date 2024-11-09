package store

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sing3demons/todoapi/logger"
	"github.com/sing3demons/todoapi/model"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type Storer interface {
	Create(*model.Todo, logger.ILogDetail) error
	List(opt FindOption, logger logger.ILogDetail) ([]model.Todo, error)
	Delete(id string, logger logger.ILogDetail) error
	FindOne(id string, logger logger.ILogDetail) (*model.Todo, error)
}

type FindOption struct {
	SearchItem  map[string]interface{}
	CommandName string
	SortItem    map[string]interface{}
	SelectItem  []string
}

type Store struct {
	sql    *gorm.DB
	mongo  *mongo.Collection
	logger logger.ILogDetail
}

type RequestLog struct {
	Body struct {
		Collection string `json:"Collection,omitempty"`
		Table      string `json:"Table,omitempty"`
		Method     string `json:"Method"`
		Query      any    `json:"Query"`
		Document   any    `json:"Document"`
		Options    any    `json:"Options"`
		Order      any    `json:"Order,omitempty"`
	} `json:"Body"`
	RawData string `json:"RawData,omitempty"`
	Data    any    `json:"Data,omitempty"`
}

func (tx *Store) Create(commandName, name, method string, data any) error {
	node := "db"
	reqLog := RequestLog{}
	reqLog.Body.Method = method
	reqLog.Body.Query = nil
	reqLog.Body.Document = data
	reqLog.Body.Options = nil
	if tx.mongo != nil {
		node = "mongo"
		jsonDocumentBytes, _ := json.Marshal(data)
		jsonDocument := strings.ReplaceAll(string(jsonDocumentBytes), `"`, "'")
		reqLog.Body.Collection = name
		reqLog.RawData = fmt.Sprintf("%s.%s(%s)", name, method, jsonDocument)

		tx.logger.AddOutput(node, commandName, reqLog).End()

		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		r, err := tx.mongo.InsertOne(ctx, data)
		if err != nil {
			tx.logger.AddError(node, commandName, "input", data, err)
			return err
		}
		reqLog.Data = r
	} else {
		node = "gorm"
		v := reflect.ValueOf(data)
		if v.Kind() == reflect.Ptr {
			v = v.Elem()
		}
		columns := make([]string, 0)
		values := make([]string, 0)
		if v.Kind() == reflect.Struct {
			for i := 0; i < v.NumField(); i++ {
				fieldName := v.Type().Field(i).Name
				fieldValue := v.Field(i).Interface()
				columns = append(columns, fieldName)
				values = append(values, fmt.Sprintf("'%v'", fieldValue))

			}
		}
		reqLog.RawData = fmt.Sprintf("insert into %s (%s) values (%s)", name, strings.Join(columns, ","), strings.Join(values, ","))
		reqLog.Body.Table = name

		tx.logger.AddOutput(node, commandName, reqLog)

		if err := tx.sql.Create(data).Error; err != nil {
			tx.logger.AddError(node, commandName, "input", data, err)
			return err
		}

		reqLog.Data = data
	}
	tx.logger.AddInput(node, commandName, reqLog.Data)
	return nil
}
