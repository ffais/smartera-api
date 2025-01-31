package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/smartcommunitylab/smartera-api/pkg/config"
	"github.com/smartcommunitylab/smartera-api/pkg/models"
)

func main() {
	config := config.ReadConfig()
	store := models.NewFileStore(config.File.Directory)
	log.Printf("url: %s, db: %s, coll: %s", config.Database.Url, config.Database.DbName, config.Database.Collection)
	// store := models.NewMongoStore(config.Database.Url, config.Database.DbName, config.Database.Collection)
	log.Printf("mongo ready")
	userHandler := NewUserHandler(store)
	home := homeHandler{}

	router := mux.NewRouter()

	router.HandleFunc("/", home.ServeHTTP)
	router.HandleFunc("/user/{id}", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/user/{id}", userHandler.GetUser).Methods("GET")
	router.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")

	http.ListenAndServe(fmt.Sprintf(": %s", config.Server.Port), router)
	log.Printf("server ready")
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("500 Internal Server Error"))
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}

type UserHandler struct {
	store userStore
}

func NewUserHandler(s userStore) *UserHandler {
	return &UserHandler{
		store: s,
	}
}

type userStore interface {
	Add(name string, user models.User) error
	Get(name string) (models.User, error)
	Update(name string, user models.User) error
}

type homeHandler struct{}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Smartera-API"))
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	id := mux.Vars(r)["id"]
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Add(id, user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	user, err := h.store.Get(id)
	if err != nil {
		if err == models.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(user)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	if err := h.store.Update(id, user); err != nil {
		if err == models.ErrNotFound {
			NotFoundHandler(w, r)
			return
		}

		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
