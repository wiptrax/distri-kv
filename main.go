package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/wiptrax/dsitributed-kv-store/config"
	"github.com/wiptrax/dsitributed-kv-store/db"
	"github.com/wiptrax/dsitributed-kv-store/web"
)

var (
	dbLocation = flag.String("db-location", "", "The path to bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("shard", "", "the name of shard for the data")
	replica    = flag.Bool("replica", false, "Whether or not run as a read-only replica")
)

func parseFlags() {
	flag.Parse()
	log.Println(*dbLocation)

	if *dbLocation == "" {
		log.Fatal("Must provide db-location")
	}

	if *shard == "" {
		log.Fatal("Must provide shard")
	}
}

func main() {
	parseFlags()

	c, err := config.Parsefile(*configFile)
	if err != nil {
		log.Fatalf("Error parsing config %q: %v", *configFile, err)
	}

	shards, err := config.ParseShards(c.Shards, *shard)
	if err != nil {
		log.Fatalf("Error parsing shards config: %v", err)
	}

	log.Printf("Shard count is %d, current shard: %d", shards.Count, shards.CurIdx)

	db, close, err := db.NewDatabase(*dbLocation, *replica)
	if err != nil {
		log.Fatalf("NewDatabse(%q) :%v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(db, shards)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)
	http.HandleFunc("/purge", srv.DeleteExtraKeysHandler)
	http.HandleFunc("/next-replication-key", srv.GetNextKeyForReplication)
	http.HandleFunc("/delete-replication-key", srv.DeleteReplicationKey)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
