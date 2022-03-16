package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(infile string) string {
	contents, err := ioutil.ReadFile(infile)
	check(err)
	return string(contents)
}

type server struct{}

func (s *server) routes() {
	http.HandleFunc("/hi", s.handleHi())
	http.HandleFunc("/inc", s.handleInc())
	http.HandleFunc("/parsejson", s.handleParseJson())
	http.HandleFunc("/get", s.handleGET())
	http.Handle("/html/", http.StripPrefix("/html/", http.FileServer(http.Dir("./htm"))))
}

func (s *server) handleHi() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("I said hi")
		fmt.Fprintf(w, "Hi from /%v\n", r.URL.Path[1:])
	}
}

func (s *server) handleInc() http.HandlerFunc {
	var i int
	return func(w http.ResponseWriter, r *http.Request) {
		prev := i
		i++
		fmt.Printf("Increment: %v -> %v\n", prev, i)
		fmt.Fprintf(w, "<p>Visited: <b>%v</b></p>\n", i)
	}
}

func (s *server) handleParseJson() http.HandlerFunc {
	type idname struct {
		Id   int
		Name string
	}
	var in idname
	json.Unmarshal([]byte(`{"id":1,"name":"test1"}`), &in)
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(in)
		fmt.Fprintf(w, "%v,%v\n", in.Id, in.Name)
	}
}

// TODO: make better
func (s *server) handleGET() http.HandlerFunc {
	idnames := map[int]string{1: "test1", 2: "test2", 3: "test3", 4: "test4", 5: "test5"}
	enc, err := json.Marshal(idnames)
	check(err)
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("%v\n", `{"id":1,"name":"test1"}`)
		// fmt.Fprintf(w, "%v\n", `{"id":1,"name":"test1"}`)
		fmt.Printf("%s\n", enc)
		fmt.Fprintf(w, "%s\n", enc)
	}
}

// TODO: post request

func main() {
	s := server{}
	s.routes()
	func() {
		serving := ":8081"
		fmt.Printf("Serving %s\n", serving)
		http.ListenAndServe(serving, nil)
	}()
}
