package web

import (
	"fmt"
	"net/http"

	"github.com/wiptrax/sitributed-kv-store/db"
)

// Server contains HTTP method handler to be used for the database
type Server struct {
	db *db.DataBase
}

// NewServer create a new server instance with HTTP handlers to get and set values
func NewServer(db *db.DataBase) *Server {
	return &Server{
		db: db,
	}
}

// GetHandler handles read request from the database
func (s *Server) GetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")

	value, err := s.db.GetKey(key)

	fmt.Fprintf(w, "Value = %v, error = %v", value, err)
}

// SetHandler handles write requests from database
func (s *Server) SetHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	key := r.Form.Get("key")
	value := r.Form.Get("value")

	// fmt.Println(key, value)
	err := s.db.SetKey(key, []byte(value))
	fmt.Fprintf(w, "error = %v", err)
}
