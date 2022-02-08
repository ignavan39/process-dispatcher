package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

var i int

type ProcessDispatcher struct {
	Status map[int]chan int
}

func (s *ProcessDispatcher) Process(id int) {
	go func() {
		for {
			select {
			default:
				time.Sleep(time.Second)
				i++
				fmt.Println(i)

			case <-s.Status[id]:
				return
			}
		}
	}()
}

func main() {

	m := ProcessDispatcher{
		Status: map[int]chan int{},
	}

	http.HandleFunc("/start", func(w http.ResponseWriter, r *http.Request) {
		if m.Status[1] == nil {
			m.Status[1] = make(chan int)
			m.Process(1)
			return
		} else {
			w.Write([]byte("Don start"))
			return
		}

	})

	http.HandleFunc("/stop", func(w http.ResponseWriter, r *http.Request) {
		if m.Status[1] == nil {
			w.Write([]byte("Dont stop is null"))
			return
		} else {
			m.Status[1] <- 1
			close(m.Status[1])
			delete(m.Status, 1)
			w.Write([]byte("Stopped"))
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
