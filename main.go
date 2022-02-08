package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type ProcessDispatcher struct {
	Status   map[int]chan int
	Counters map[int]int
}

type HandleProcessResponse struct {
	Id int `json:"id"`
}

func (s *ProcessDispatcher) Process(id int) {
	go func() {
		for {
			select {
			default:
				time.Sleep(time.Second)
				s.Counters[id] = s.Counters[id] + 1
				fmt.Println(s.Counters[id])

			case <-s.Status[id]:
				return
			}
		}
	}()
}

func main() {

	m := ProcessDispatcher{
		Status:   map[int]chan int{},
		Counters: map[int]int{},
	}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		var payload HandleProcessResponse
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if m.Status[payload.Id] == nil {
			m.Status[payload.Id] = make(chan int)
			m.Process(payload.Id)

			w.WriteHeader(http.StatusOK)
			return
		} else {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("process has already been started"))
			return
		}

	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		var payload HandleProcessResponse
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {

			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		if m.Status[payload.Id] == nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("process has already been stopped"))
			return
		} else {
			m.Status[payload.Id] <- 1
			close(m.Status[payload.Id])
			delete(m.Status, payload.Id)

			w.WriteHeader(http.StatusOK)
			return
		}
	})

	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}
