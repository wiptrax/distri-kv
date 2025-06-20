package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"github.com/wiptrax/dsitributed-kv-store/config"
	"github.com/wiptrax/dsitributed-kv-store/db"
	"github.com/wiptrax/dsitributed-kv-store/web"
)

var (
	dbLocation = flag.String("db-location", "", "The path to bolt db database")
	httpAddr   = flag.String("http-addr", "127.0.0.1:8080", "HTTP host and port")
	configFile = flag.String("config-file", "sharding.toml", "Config file for static sharding")
	shard      = flag.String("shard", "", "the name of shard for the data")
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

	var c config.Config
	if _, err := toml.DecodeFile(*configFile, &c); err != nil {
		log.Fatalf("toml.Decodefile(%q): %v", *configFile, err)
	}

	shardCount := len(c.Shards)
	var shardIdx int = -1
	var addrs = make(map[int]string)


	for _, s := range c.Shards {
		addrs[s.Idx] = s.Address

		if s.Name == *shard {
			shardIdx = s.Idx
		}
	}

	if shardIdx < 0 {
		log.Fatalf("shard %q was not found", *shard)
	}

	log.Printf("Shard caount is %d, current shard is %d", shardCount, shardIdx)

	db, close, err := db.NewDatabase(*dbLocation)
	if err != nil {
		log.Fatalf("NewDatabse(%q) :%v", *dbLocation, err)
	}
	defer close()

	srv := web.NewServer(db, shardIdx, shardCount, addrs)

	http.HandleFunc("/get", srv.GetHandler)
	http.HandleFunc("/set", srv.SetHandler)

	log.Fatal(http.ListenAndServe(*httpAddr, nil))
}
