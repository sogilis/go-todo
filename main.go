package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"flag"
	"path"
	"log"
	"fmt"
	"os"
)

type TodoItem struct {
	Name string
	Completed bool
}

var items []TodoItem = []TodoItem{}

func handleRoot(w http.ResponseWriter, req *http.Request) {
  // The "/" pattern matches everything, so we need to check
  // that we're at the root here.
  if req.URL.Path != "/" {
      http.NotFound(w, req)
      return
  }

  fmt.Fprintf(w, "Welcome to the home page!")
}

func handleWebError(w http.ResponseWriter, e error, status int) {
	w.WriteHeader(status)
	fmt.Fprintf(w, "Error: %v", e)
}

func listItems(w http.ResponseWriter, req *http.Request) {

	switch(req.Method) {
	case "POST":
		data, err := ioutil.ReadAll(req.Body)

		if err != nil {
			handleWebError(w, err, http.StatusInternalServerError)
			return
		}

		items = append(items, TodoItem{
			Name: string(data),
			Completed: false,
		})
		
		err = save()
		if err != nil {
			handleWebError(w, err, http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)

	case "HEAD": fallthrough
	case "GET":
		res, err := json.Marshal(items)

		if err != nil {
			handleWebError(w, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json; encoding=\"utf-8\"")
		w.Write(res)
	}
}

func singleItem(w http.ResponseWriter, req *http.Request) {
	lastItem := path.Base(req.URL.Path)
	n, err := strconv.Atoi(lastItem)

	if err != nil {
		handleWebError(w, err, http.StatusBadRequest)
	} else if n < 0 || n >= len(items) {
		handleWebError(w, fmt.Errorf("Index out of range 0..%d", len(items)-1), http.StatusBadRequest)
		return
	}

	if req.Method == "DELETE" {
		
		items = append(items[:n], items[n+1:]...)
		
		err = save()
		if err != nil {
			handleWebError(w, err, http.StatusInternalServerError)
			return
		}
		
		w.WriteHeader(http.StatusNoContent)
		
	} else {	

		res, err := json.Marshal(items[n])
	
		if err != nil {
			handleWebError(w, err, http.StatusInternalServerError)
			return
		}
	
		w.Header().Set("Content-Type", "application/json; encoding=\"utf-8\"")
		w.Write(res)
	
	}
}

func save() error {
	f, err := os.Create("items.json")
	if err != nil {
		return err
	}
	defer f.Close()
	
	res, err := json.Marshal(items)
	if err != nil {
		return err
	}
	
	_, err = f.Write(res)
	if err != nil {
		return err
	}
	
	return nil
}

func open() error {
	
	f, err := os.Open("items.json")
	if err != nil {
		return err
	}
	defer f.Close()
	
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	
	err = json.Unmarshal(data, &items)
	if err != nil {
		return err
	}
	
	return nil
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("/list", listItems)
	mux.HandleFunc("/list/", singleItem)
}

func main() {
	arg_port := flag.Int("port", 8080, "HTTP Port")
	flag.Parse()
	err := open()
	if err != nil {
		log.Fatal(err)
	}	
	log.Printf("Starting TODO server on port %d\n", *arg_port)
	mux := http.NewServeMux()
	registerRoutes(mux)
	srv := http.Server{
		Addr: fmt.Sprintf(":%d", *arg_port),
		Handler: mux,
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}


