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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (tx *Store) List(commandName, name string, opt FindOption, data any) (interface{}, error) {
	node := "db"
	reqLog := RequestLog{}
	reqLog.Body.Method = "find"

	reqLog.Body.Document = nil
	reqLog.Body.Options = nil

	if tx.mongo != nil {
		r, err := tx.listFromMongo(commandName, name, opt, data, reqLog)
		if err != nil {
			tx.logger.AddError(node, commandName, "input", nil, err)
			return nil, err
		}
		return r, nil

	} else if tx.sql != nil {
		return tx.listFromSQL(commandName, name, opt, data, reqLog)
	}

	tx.logger.AddInput(node, commandName, data)
	return data, nil
}

func (tx *Store) listFromMongo(commandName, name string, opt FindOption, data any, reqLog RequestLog) (interface{}, error) {
	node := "mongo"
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	filter := buildMongoFilter(opt)
	opts := buildMongoFindOptions(opt)

	reqLog.Body.Query = filter
	reqLog.Body.Options = opts
	reqLog.Body.Collection = name

	reqLog.RawData = buildMongoRawData(name, filter, opts)

	tx.logger.AddOutput(node, commandName, reqLog).End()

	cur, err := tx.mongo.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	if err := cur.All(ctx, &data); err != nil {
		return nil, err
	}

	tx.logger.AddInput(node, commandName, data)
	return data, nil
}

func buildMongoFilter(opt FindOption) bson.D {
	filter := bson.D{{Key: "deleted_at", Value: primitive.Null{}}}
	if opt.SearchItem != nil {
		for k, v := range opt.SearchItem {
			filter = append(filter, bson.E{Key: k, Value: v})
		}
	}
	return filter
}

func buildMongoFindOptions(opt FindOption) *options.FindOptions {
	opts := &options.FindOptions{Sort: bson.D{{Key: "created_at", Value: -1}}}
	if opt.SelectItem != nil {
		projection := bson.D{}
		for _, v := range opt.SelectItem {
			projection = append(projection, bson.E{Key: v, Value: 1})
		}
		opts.Projection = projection
	}
	if opt.SortItem != nil {
		for k, v := range opt.SortItem {
			if v == "asc" {
				opts.Sort = bson.D{{Key: k, Value: 1}}
			} else {
				opts.Sort = bson.D{{Key: k, Value: -1}}
			}
		}
	}
	return opts
}

func buildMongoRawData(name string, filter bson.D, opts *options.FindOptions) string {
	rawData := fmt.Sprintf("%s.find(%s, {projection,sort})", name, ConvertDToJSON(filter))
	if opts.Sort != nil {
		rawData = strings.Replace(rawData, "sort", fmt.Sprintf("sort:%s", ConvertDToJSON(opts.Sort.(bson.D))), 1)
	} else {
		rawData = strings.Replace(rawData, "sort", "", 1)
	}
	if opts.Projection != nil {
		rawData = strings.Replace(rawData, "projection", fmt.Sprintf("projection:%s", ConvertDToJSON(opts.Projection.(bson.D))), 1)
	} else {
		rawData = strings.Replace(rawData, "projection,", "", 1)
	}

	if opts.Projection == nil && opts.Sort == nil {
		rawData = strings.Replace(rawData, "{", "", 1)
		rawData = strings.Replace(rawData, "}", "", 1)
		rawData = strings.Replace(rawData, ", ", "", 1)
	}
	fmt.Println("RawData=========================", rawData)
	return rawData
}

func (tx *Store) listFromSQL(commandName, name string, opt FindOption, data any, reqLog RequestLog) (interface{}, error) {
	node := "gorm"
	var conds []interface{}
	var cond []string
	if opt.SearchItem != nil {
		for k, v := range opt.SearchItem {
			conds = append(conds, fmt.Sprintf("%s = ?", k), v)
			cond = append(cond, fmt.Sprintf("%s = '%v'", k, v))
		}
	}

	var order []string
	if opt.SortItem != nil {
		for k, v := range opt.SortItem {
			order = append(order, fmt.Sprintf("%s %s", k, v))
		}
	}

	selectTodo := []string{"*"}
	if opt.SelectItem != nil {
		selectTodo = opt.SelectItem
	}

	reqLog.Body.Order = order
	reqLog.Body.Query = conds
	reqLog.Body.Document = selectTodo
	rawData := fmt.Sprintf("SELECT %s FROM %s", strings.Join(selectTodo, ","), name)
	if cond != nil {
		rawData = fmt.Sprintf("%s WHERE %s", rawData, strings.Join(cond, " and "))
	}

	if order != nil {
		rawData = fmt.Sprintf("%s order by %s", rawData, strings.Join(order, ","))
	}
	reqLog.RawData = rawData

	tx.logger.AddOutput(node, commandName, reqLog).End()
	fmt.Println("List=========================", reqLog.RawData)

	r := tx.sql.Select(selectTodo).Order(strings.Join(order, ",")).Find(&data, conds...)
	if err := r.Error; err != nil {
		tx.logger.AddError(node, commandName, "input", nil, err)
		return nil, err
	}

	tx.logger.AddInput(node, commandName, data)
	return data, nil
}

func ConvertDToJSON(d bson.D) string {
	// Marshal primitive.D to JSON-compatible byte slice
	data := convertDToMap(d)
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error converting map to JSON:", err)
		return ""
	}
	return string(jsonData)
}

func convertDToMap(d bson.D) map[string]interface{} {
	result := make(map[string]interface{})
	for _, elem := range d {
		result[elem.Key] = elem.Value
	}
	return result
}
