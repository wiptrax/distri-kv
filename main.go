package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wiptrax/sitributed-kv-store/db"
	"github.com/wiptrax/sitributed-kv-store/web"
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

	srv := web.NewServer(db)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
