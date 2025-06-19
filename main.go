package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/wiptrax/sitributed-kv-store/db"
)

var (
	dbLocation = flag.String("db-location", "", "The path to bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
)

func parseFlags() {
	flag.Parse()
	log.Println(*dbLocation)

	if *dbLocation == "" {
		log.Fatal("Must provide db-location")
	}
}

func main() {
	parseFlags()

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatal("NewDatabse(%q) :%v", *dbLocation, err)
	}
	defer close()

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.Form.Get("key")

		value, err := db.GetKey(key)

		fmt.Fprintf(w, "Value = %v, error = %v", value, err)
	})
	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		key := r.Form.Get("key")
		value := r.Form.Get("value")

		// fmt.Println(key, value)
		err := db.SetKey(key, []byte(value))
		fmt.Fprintf(w, "error = %v", err)
	})

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
