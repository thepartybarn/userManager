package main

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
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

	myMux.Handle("/api/getToken", getTokenHandler, "NONE")
	myMux.Handle("/api/closeToken", closeTokenHandler, "NONE")
	myMux.Handle("/api/getUserInfoForToken", getUserInfoForTokenHandler, "NONE")

	myMux.Handle("/api/createUser", createUserHandler, "CREATE_USER")
	myMux.Handle("/api/deleteUser", deleteUserHandler, "DELETE_USER")

	myMux.Handle("/api/addPermissionToUser", addPermissionToUserHandler, "ADD_PERMISSION_TO_USER")
	myMux.Handle("/api/removePermissionFromUser", removePermissionFromUserHandler, "Remove_PERMISSION_FROM_USER")

	log.Trace("Opening HTTP Server")
	err = http.ListenAndServe(":80", myMux)
	if err != nil {
		log.Panic(err)
	}
}

func addPermissionToUserHandler(context *Context) {

}
func removePermissionFromUserHandler(context *Context) {

}
func createUserHandler(context *Context) {
	username := context.Parameters.Get("username")

	context.ReturnData["username"] = "test!" + username

	context.returnJson(http.StatusOK)
}
func deleteUserHandler(context *Context) {

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

func getUserInfoForTokenHandler(context *Context) {
	context.ReturnData = getUserInfoForToken(context.Parameters.Get("token"))
	context.returnJson(http.StatusOK)
}
func getTokenHandler(context *Context) {
	username := context.Parameters.Get("username")
	password := context.Parameters.Get("password")
	cardID := context.Parameters.Get("cardID")
	log.Tracef("getTokenHandler (Username: %s, Password: %s, CardID: %s)", username, password, cardID)

	context.ReturnData["token"] = getToken(username, password, cardID)
	context.returnJson(http.StatusOK)
}
func closeTokenHandler(context *Context) {
	closeToken(context.Parameters.Get("token"))
	context.returnJson(http.StatusOK)
}
