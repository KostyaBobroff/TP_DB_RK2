package service

import (	
	"fmt"
	"net/http"
	"io/ioutil"
	"../database"
	"encoding/json"
	//"github.com/bozaro/tech-db-forum/generated/client/operations"
	// "github.com/go-openapi/strfmt"
	"github.com/bozaro/tech-db-forum/generated/models"
	"github.com/gorilla/mux"
	"github.com/go-openapi/swag"
)

// /forum/create Создание форума
func CreateForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateForum")
	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	forum := &models.Forum{}
	err = json.Unmarshal(body, &forum)

	// reg := strfmt.NewFormats()
	// err = forum.Validate(reg)
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.CreateForumDB(forum)

	resp, err1 := result.MarshalBinary()
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err1)
	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte("Hello1 "))
	case database.ForumIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}

// /forum/{slug}/details Получение информации о форуме
func GetForum(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetForum")
	params := mux.Vars(r)
	slug := params["slug"]

	result, err := database.GetForumDB(slug)

	resp, err1 := result.MarshalBinary()
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err1)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /forum/{slug}/create Создание ветки
func CreateForumThread(w http.ResponseWriter, r *http.Request) {
	fmt.Println("CreateForumThread")
	params := mux.Vars(r)
	slug := params["slug"]

	body, err := ioutil.ReadAll(r.Body)	
	defer r.Body.Close()
	if err != nil {
		fmt.Println(err)
		return
	}	
	thread := &models.Thread{}
	err = json.Unmarshal(body, &thread)
	thread.Slug = slug
	thread.Forum = slug // иначе не знаю как

	//err = forum.Validate()
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := database.CreateForumThreadDB(thread)

	resp, _ := result.MarshalBinary()
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 201, resp)
	case database.UserNotFound:
		makeResponse(w, 404, []byte("Hello1 "))
	case database.ForumIsExist:
		makeResponse(w, 409, resp)
	default:		
		makeResponse(w, 500, []byte("Hello2 "))
	}
}


// /forum/{slug}/threads Список ветвей обсужления форума
func GetForumThreads(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetForumThreads")
	params := mux.Vars(r)
	slug := params["slug"]
	queryParams := r.URL.Query()
	var limit, since, desc string
	if limit = queryParams.Get("limit"); limit != "" {
		limit = "1";
	}
	if since = queryParams.Get("since"); limit != "" {
		since = "";
	}
	if desc = queryParams.Get("desc"); limit != ""{
		desc = "false";
	}
	fmt.Println(limit, since, desc)

	result, err := database.GetForumThreadsDB(slug, limit, since, desc)
	fmt.Println(result)
	fmt.Println(err)
	
	// resp, _ := result.MarshalBinary()
	resp, _ := swag.WriteJSON(result)
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}

// /forum/{slug}/users Пользователи данного форума
func GetForumUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetForumUsers")
	params := mux.Vars(r)
	slug := params["slug"]
	queryParams := r.URL.Query()
	var limit, since, desc string
	if limit = queryParams.Get("limit"); limit != "" {
		limit = "1";
	}
	if since = queryParams.Get("since"); limit != "" {
		since = "";
	}
	if desc = queryParams.Get("desc"); limit != ""{
		desc = "false";
	}
	fmt.Println(limit, since, desc)

	result, err := database.GetForumUsersDB(slug, limit, since, desc)
	fmt.Println(result)
	fmt.Println(err)
	
	// resp, _ := result.MarshalBinary()
	resp, _ := swag.WriteJSON(result)
	fmt.Println("DB result")
	fmt.Println(string(resp))
	fmt.Println(err)
	switch err {
	case nil:
		makeResponse(w, 200, resp)
	case database.ForumNotFound:
		makeResponse(w, 404, []byte("Can't find user with id #42\n"))
	default:		
		makeResponse(w, 500, []byte("Hello here"))
	}
}