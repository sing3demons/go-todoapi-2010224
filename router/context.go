package router

import "github.com/sing3demons/todoapi/logger"

type IContext interface {
	Bind(interface{}) error
	JSON(int, interface{})
	Log(name string) logger.ILogDetail
	Get(string) interface{}
	TransactionID() string
	Param(string) string
	Query(string) string
}
