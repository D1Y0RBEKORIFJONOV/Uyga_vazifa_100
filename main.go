package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"uyga_vazifa3/db"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

var queries *db.Queries

func main() {
	database, err := sql.Open("postgres", "postgresql://postgres:+_+diyor2005+_+@localhost:5432/task?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	queries = db.New(database)

	r := mux.NewRouter()

	// Routelarni ro'yxatdan o'tkazish
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks", listTasks).Methods("GET")
	r.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
	r.HandleFunc("/authors", createAuthor).Methods("POST")
	r.HandleFunc("/authors", listAuthors).Methods("GET")
	r.HandleFunc("/authors/{id}", getAuthor).Methods("GET")
	r.HandleFunc("/tasks/{task_id}/authors/{author_id}", assignAuthorToTask).Methods("POST")
	log.Println("listen server: 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task db.CreateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTask, err := queries.CreateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdTask)
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := queries.ListTasks(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(tasks)
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	task, err := queries.GetTask(r.Context(), db.GetTaskParams{
		int32(id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(task)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	var task db.UpdateTaskParams
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.ID = int32(id)

	updatedTask, err := queries.UpdateTask(r.Context(), task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	err := queries.DeleteTask(r.Context(), int32(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func createAuthor(w http.ResponseWriter, r *http.Request) {
	var author db.CreateAuthorParams
	if err := json.NewDecoder(r.Body).Decode(&author); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdAuthor, err := queries.CreateAuthor(r.Context(), author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdAuthor)
}

func listAuthors(w http.ResponseWriter, r *http.Request) {
	authors, err := queries.ListAuthors(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(authors)
}

func getAuthor(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])
	author, err := queries.GetAuthor(r.Context(), db.GetAuthorParams{
		int32(id),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(author)
}

func assignAuthorToTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	taskID, _ := strconv.Atoi(vars["task_id"])
	authorID, _ := strconv.Atoi(vars["author_id"])

	err := queries.AssignAuthorToTask(r.Context(), db.AssignAuthorToTaskParams{
		TaskID:   int32(taskID),
		AuthorID: int32(authorID),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
