package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx"
)

type server struct {
	db *pgx.ConnPool
}

func newServer(pool *pgx.ConnPool) *server {
	return &server{
		db: pool,
	}
}

func (s *server) handelList() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", 405)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var req GetStationsListRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Can't decode input request: %s", err.Error())
			http.Error(w, "Can't decode input request", 500)
			return
		}
		if req.Train == "" {
			http.Error(w, "train can't be empty", 400)
			return
		}
		st, err := getStationList(s.db, req.Train)
		if err != nil {
			log.Printf("Can't get station list: %s", err.Error())
			http.Error(w, "Can't get station list", 500)
			return
		}

		if st == nil || len(st.Stations) == 0 {
			http.Error(w, "Not Found", 404)
			return
		}
		result, err := json.Marshal(st)
		if err != nil {
			log.Printf("Can't encode response: %s", err.Error())
			http.Error(w, "Can't encode response", 500)
			return
		}
		fmt.Fprintf(w, string(result))
	}
}

func (s *server) handelOrder() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", 405)
			return
		}
		decoder := json.NewDecoder(r.Body)
		var req PostOrderRequest
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Can't decode input request: %s", err.Error())
			http.Error(w, "Can't decode input request", 500)
			return
		}
		if req.Train == "" {
			http.Error(w, "train can't be empty", 400)
			return
		}
		resp, err := createOrder(s.db, &req)
		if err != nil {
			log.Printf("Can't create new order: %s", err.Error())
			http.Error(w, "Can't create new order", 500)
			return
		}

		result, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Can't encode response: %s", err.Error())
			http.Error(w, "Can't encode response", 500)
			return
		}
		fmt.Fprintf(w, string(result))
	}
}

func (s *server) handelTasks() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Invalid request method.", 405)
			return
		}

		tasks, err := getTaskList(s.db)
		if err != nil {
			log.Printf("Can't get task list: %s", err.Error())
			http.Error(w, "Can't get task list", 500)
			return
		}

		if tasks == nil || len(tasks.Tasks) == 0 {
			http.Error(w, "Not Found", 404)
			return
		}
		result, err := json.Marshal(tasks)
		if err != nil {
			log.Printf("Can't encode response: %s", err.Error())
			http.Error(w, "Can't encode response", 500)
			return
		}
		fmt.Fprintf(w, string(result))
	}
}
