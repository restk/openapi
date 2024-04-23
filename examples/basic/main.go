package main

import (
	"fmt"
	"net/http"

	"github.com/restk/openapi"
)

func main() {
	type Error struct {
		Message string `json:"message"`
	}
	type Success struct {
		Message string `json:"message"`
	}

	openAPI := openapi.New("My API", "v1.0.0")
	openAPI.Description("My API is a great API")
	openAPI.Server().URL("https://myapi.com").Description("My API URL")

	// define security
	openAPI.BearerAuth()

	// define security
	openAPI.OAuth2().Implicit().AuthorizationURL("https://api.example.com/oauth2/authorize").Scopes(map[string]string{
		"read_users":  "read users",
		"write_users": "write to users",
	})

	// apply security (globally)
	openAPI.Security("OAuth2", []string{"read_users", "write_users"})

	// registering a URL (operation)
	type LoginRequest struct {
		Username string `json:"username" doc:"username to login with" maxLength:"25"`
		Password string `json:"password" doc:"password to login with" maxLength:"50"`
	}

	login := openAPI.Register(&openapi.Operation{
		Method: "POST",
		Path:   "/login",
	})

	login.Request().Body(&LoginRequest{})
	login.Response(http.StatusOK).ContentType("text/plain").Body(openapi.StringType)
	login.Response(http.StatusOK).Body(Success{})

	// registering a URL (operation)

	type User struct {
		ID   string `json:"id" doc:"ID of user" example:"2"`
		Name string `json:"name" doc:"Name of user" example:"joe"`
		Age  int    `json:"age" doc:"Age of user" example:"32"`
	}

	getUser := openAPI.Register(&openapi.Operation{
		Method: "GET",
		Path:   "/login/{userId}",
	})

	getUser.Request().PathParam("userId", openapi.IntType).Required(true)
	getUser.Request().QueryParam("age", &openapi.IntType).Example("12").Description("Age of user") // pointers make the query param optional
	getUser.Request().QueryParam("email", openapi.StringType).Example("joe@gmail.com").Description("Email of user")

	getUser.Response(http.StatusOK).Body(User{}).Example(`{id: 3, name: "joe", age: 5}`) // override example from User struct
	getUser.Response(http.StatusForbidden).Body(Error{})                                 // default content type is application/json

	bytes, err := openAPI.OpenAPI().YAML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))

	// serve under /docs using scalar (visit http://localhost:8111/docs)
	scalar := openapi.Scalar(openAPI.OpenAPI(), map[string]any{
		"theme": "purple",
	})

	fmt.Println("Serving docs at http://localhost:8111/docs")

	docs := func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(scalar)
	}

	http.HandleFunc("/docs", docs)
	http.ListenAndServe(":8111", nil)
}
