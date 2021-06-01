package svc

import (
	"github.com/tal-tech/go-zero/core/logx"
	"github.com/tal-tech/go-zero/rest"
	"greet/internal/config"
	"net/http"
)

type ServiceContext struct {
	Config           config.Config
	GreetMiddleware1 rest.Middleware
	GreetMiddleware2 rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:           c,
		GreetMiddleware1: greetMiddleware1,
		GreetMiddleware2: greetMiddleware2,
	}
}

func greetMiddleware1(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("greetMiddleware1 request ... ")
		next(w, r)
		logx.Info("greetMiddleware1 response ...")
	}
}

func greetMiddleware2(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logx.Info("greetMiddleware2 request ... ")
		next(w, r)
		logx.Info("greetMiddleware2 response ... ")
	}
}
