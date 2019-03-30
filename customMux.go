package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"
)

type Handler func(*Context)
type Route struct {
	Pattern    *regexp.Regexp
	Handler    Handler
	Permission string
}
type CustomMux struct {
	Routes       []Route
	DefaultRoute Handler
}

func (customMux *CustomMux) Handle(pattern string, handler Handler, permission string) {
	regexp := regexp.MustCompile(pattern)
	route := Route{Pattern: regexp, Handler: handler, Permission: permission}

	customMux.Routes = append(customMux.Routes, route)
}

func (customMux *CustomMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	context := &Context{Request: r, ResponseWriter: w, ReturnData: make(map[string]interface{})}
	for _, route := range customMux.Routes {
		log.Debug(route)
		if matches := route.Pattern.FindStringSubmatch(context.URL.Path); len(matches) > 0 {
			context.Parameters = context.URL.Query()
			token := Token(context.Parameters.Get("token"))
			if (route.Permission == "NONE") || (token != "" && tokenHasPermission(token, route.Permission)) {
				route.Handler(context)
			}
			return
		}
	}
	customMux.DefaultRoute(context)
}

type Context struct {
	http.ResponseWriter
	*http.Request
	Parameters url.Values
	ReturnData map[string]interface{}
}

func (context *Context) returnJson(code int) {
	context.ResponseWriter.Header().Set("Content-Type", "application/json")
	context.WriteHeader(code)

	outputDataBytes, err := json.Marshal(context.ReturnData)
	if err != nil {
		log.Error(err)
	}
	context.ResponseWriter.Write(outputDataBytes)
}
