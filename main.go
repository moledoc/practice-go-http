package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	_ "strconv"
	"strings"
	"sync"
)

// NOTE: no proper error handling/redirecting done!

func check(w http.ResponseWriter, err error) {
	if err != nil {
		panic(err)
	}
}

func ReadFile(infile string) string {
	contents, err := ioutil.ReadFile(infile)
	if err != nil {
		panic(err)
	}
	return string(contents)
}

type server struct{}

func (s *server) routes() {
	http.HandleFunc("/hi", s.handleHi())
	http.HandleFunc("/inc", s.handleInc())
	http.HandleFunc("/new/idname", s.handleNewIdname())
	http.HandleFunc("/get", s.handleGET())
	http.HandleFunc("/post", s.handlePOST())
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

func (s *server) handleNewIdname() http.HandlerFunc {
	type idname struct {
		Id   string
		Name string
	}
	parseArgs := func(r *http.Request) (map[string]string, error) {
		upArgs := strings.Split(r.URL.RawQuery, "&")
		args := make(map[string]string)
		for _, unparsed := range upArgs {
			kv := strings.Split(unparsed, "=")
			k, v := kv[0], kv[1]
			// k, err := strconv.Atoi(kv[0])
			// if err != nil {
			// return map[string]string{}, err
			// }
			_, ok := args[k]
			if ok {
				return map[string]string{}, errors.New("Redefinition of parameter")
			}
			args[k] = v
		}
		return args, nil
	}
	// var in idname
	// json.Unmarshal([]byte(`{"id":1,"name":"test1"}`), &in)
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(in)
		// fmt.Fprintf(w, "%v,%v\n", in.Id, in.Name)
		args, err := parseArgs(r)
		check(w, err)
		j, err := json.Marshal(args)
		check(w, err)
		fmt.Fprintf(w, string(j))
	}
}

func (s *server) handleGET() http.HandlerFunc {
	idnames := map[int]string{1: "test1", 2: "test2", 3: "test3", 4: "test4", 5: "test5"}
	var (
		header map[string][]string
		init   sync.Once
	)
	return func(w http.ResponseWriter, r *http.Request) {
		enc, err := json.Marshal(idnames)
		check(w, err)
		init.Do(func() { // https://pace.dev/blog/2018/05/09/how-I-write-http-services-after-eight-years.html
			header = w.Header()
			header["Content-Type"] = []string{"application/json"}
		})

		_, err = w.Write(enc)
		check(w, err)
		// fmt.Printf("%v\n", `{"id":1,"name":"test1"}`)
		// fmt.Fprintf(w, "%v\n", `{"id":1,"name":"test1"}`)
		fmt.Printf("%s\n", enc)
		fmt.Println(*r)
		// fmt.Fprintf(w, "%s\n", enc)
	}
}

func (s *server) handlePOST() http.HandlerFunc {
	type idname struct {
		Id   int
		Name string
	}
	var b idname
	return func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println(*r)
		fmt.Println((*r).Method)
		fmt.Println((*r).Header)
		check(w, json.NewDecoder(r.Body).Decode(&b))
		fmt.Println(b)
		fmt.Fprintf(w, "%v: Data received\n", http.StatusOK)
	}
}

func main() {
	s := server{}
	s.routes()
	func() {
		serving := ":8081"
		fmt.Printf("Serving %s\n", serving)
		http.ListenAndServe(serving, nil)
	}()
}
