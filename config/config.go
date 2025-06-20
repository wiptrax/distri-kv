package config

// Each Shard have uniques set of keys
type Shard struct {
	Name string
	Idx  int
	Address string
}

// Config is describes the sharding config
type Config struct {
	Shards []Shard
}
