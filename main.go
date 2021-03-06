package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/context"
	"github.com/justinas/alice"
	"gopkg.in/mgo.v2"
)

func recoverHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %+v", err)
				http.Error(w, http.StatusText(500), 500)
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func loggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		next.ServeHTTP(w, r)
		t2 := time.Now()
		log.Printf("[%s] %q %v\n", r.Method, r.URL.String(), t2.Sub(t1))
	}

	return http.HandlerFunc(fn)
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the about page.")
	fmt.Fprintf(w, "Params: %v", context.GetAll(r))
}

func main() {
	log.Printf("Connecting to database...")
	session, err := mgo.Dial("localhost")
	if err != nil {
		log.Printf("Failed to connect to server: %v\n", err)
		return
	}
	defer session.Close()
	common := alice.New(context.ClearHandler, loggingHandler, recoverHandler)

	router := NewRouter()
	router.Get("/about/:id", common.ThenFunc(aboutHandler))

	router.ListenAndServe(":8080")
}
