package web

import (
	"fmt"
	"io"
	"net/http"

	"github.com/wiptrax/dsitributed-kv-store/config"
	"github.com/wiptrax/dsitributed-kv-store/db"
)

// Server contains HTTP method handler to be used for the database
type Server struct {
	db     *db.DataBase
	shards *config.Shards
}

// NewServer create a new server instance with HTTP handlers to get and set values
func NewServer(db *db.DataBase, s *config.Shards) *Server {
	return &Server{
		db:     db,
		shards: s,
	}
}

func (s *Server) redirect(shard int, w http.ResponseWriter, r *http.Request) {
	url := "http://" + s.shards.Addrs[shard] + r.RequestURI
	fmt.Fprintf(w, "redirecting from shard %d to shard %d (%q)\n", s.shards.CurIdx, shard, url)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error redirecting the request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	io.Copy(w, resp.Body)
}

// GetHandler handles read request from the database
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	shard := s.shards.Index(key)
	
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}
	value, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Shard = %d, current shard = %d, addr = %q, Value = %q, error = %v", shard, s.shards.CurIdx, s.shards.Addrs[shard], value, err)
}

// SetHandler handles write requests from database
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	shard := s.shards.Index(key)
	if shard != s.shards.CurIdx {
		s.redirect(shard, w, r)
		return
	}

	// fmt.Println(key, value)
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "Error = %v, shardIdx = %d, current shard = %d", err, shard, s.shards.CurIdx)
}

// DeleteExtraKeysHandler delete keys that doesn't belong to the curren t shards
func (s *Server) DeleteExtraKeysHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Error = %v", s.db.DeleteExtraKeys(func(key string) bool {
		return s.shards.Index(key) != s.shards.CurIdx
	}))
}
