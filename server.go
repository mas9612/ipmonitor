package ipmonitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	// for gorm to use sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// NewHTTPHandler returns *mux.Router
func NewHTTPHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hosts", hostsHandler).Methods("GET")
	r.HandleFunc("/hosts/{id}", hostHandler).Methods("GET")

	return r
}

func hostsHandler(w http.ResponseWriter, r *http.Request) {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Println("[ERROR] /hosts: failed to open DB", err)
		return
	}
	defer db.Close()

	var hosts []HostModel
	err = db.Find(&hosts).Error
	if err != nil {
		log.Println("[ERROR] /hosts: failed to get host records:", err)
		return
	}

	switch r.Method {
	case http.MethodGet:
		res := make([]Host, len(hosts))
		for i, host := range hosts {
			res[i].Address = host.Address
			res[i].Hostname = host.Hostname
			res[i].Description = host.Description
		}
		replyJSON(w, http.StatusOK, HostsResponse{
			Count: len(res),
			Hosts: res,
		})
		return
	}
	replyError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s is not allowed in this URL", r.Method))
}

func hostHandler(w http.ResponseWriter, r *http.Request) {

}

func replyJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
	}
}

func replyError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	encoder := json.NewEncoder(w)
	err := encoder.Encode(ErrorResponse{
		Status:  status,
		Message: message,
	})
	if err != nil {
		log.Println("Failed to reply error")
	}
}
