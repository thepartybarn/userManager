package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"regexp"

	"github.com/sirupsen/logrus"
)

var (
	_buildDate    string
	_buildVersion string
	log           = logrus.New()
)

func main() {
	var err error
	log.SetLevel(logrus.TraceLevel)
	log.Printf("---------- Program Started %v (%v) ----------", _buildVersion, _buildDate)

	setupDataBase()

	//TODO This is temporary
	seedDatabase()

	myMux := &CustomMux{
		DefaultRoute: func(context *Context) {
			context.returnJson(http.StatusNotFound)
		},
	}

	myMux.Handle("/api/getSecretKey", getSecretKeyHandler, NONE)
	myMux.Handle("/api/closeSecretKey", closeSecretKeyHandler, NONE)
	myMux.Handle("/api/getUserInfoForSecretKey", getUserInfoForSecretKeyHandler, NONE)
	myMux.Handle("/api/createUser", createUserHandler, CREATE_USER)

	//http.HandleFunc("/api/addUser", addUserHandler)
	//http.HandleFunc("/api/addFriend", addFriendHandler)

	log.Trace("Opening HTTP Server")
	err = http.ListenAndServe(":80", myMux)
	if err != nil {
		log.Panic(err)
	}
}

func createUserHandler(context *Context) {
	username := context.Parameters.Get("username")

	context.ReturnData["username"] = "test!" + username

	context.returnJson(http.StatusOK)
}
func addFriendHandler(context *Context) {
	context.ResponseWriter.Header().Set("Content-Type", "application/json")
	context.ResponseWriter.WriteHeader(http.StatusOK)
	outputData := make(map[string]string)
	inputData := context.Request.URL.Query()

	username := inputData.Get("username")
	outputData["username"] = username

	outputDataBytes, err := json.Marshal(outputData)
	if err != nil {
		log.Error(err)
	}
	context.ResponseWriter.Write(outputDataBytes)
}
func removeFriendHandler(context *Context) {

}
func listFriendHandler(context *Context) {

}

func getUserInfoForSecretKeyHandler(context *Context) {
	context.ReturnData = getUserInfoForSecretKey(context.Parameters.Get("secretKey"))
	context.returnJson(http.StatusOK)
}
func getSecretKeyHandler(context *Context) {
	username := context.Parameters.Get("username")
	password := context.Parameters.Get("password")
	cardID := context.Parameters.Get("cardID")

	context.ReturnData["secretKey"] = getSecretKey(username, password, cardID)
	context.returnJson(http.StatusOK)
}
func closeSecretKeyHandler(context *Context) {
	closeSecretKey(context.Parameters.Get("secretKey"))
	context.returnJson(http.StatusOK)
}

type Handler func(*Context)
type Route struct {
	Pattern    *regexp.Regexp
	Handler    Handler
	Permission Permission
}
type CustomMux struct {
	Routes       []Route
	DefaultRoute Handler
}

func (customMux *CustomMux) Handle(pattern string, handler Handler, permission Permission) {
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
			secretKey := SecretKey(context.Parameters.Get("secretKey"))
			if (route.Permission == NONE) || (secretKey != "" && secretKeyHasPermission(secretKey, route.Permission)) {
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
