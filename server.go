package ipmonitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// NewHTTPHandler returns *mux.Router
func NewHTTPHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hosts", hostsHandler).Methods("GET")
	r.HandleFunc("/hosts/{id}", hostHandler).Methods("GET")

	return r
}

func hostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hosts := []Host{
			{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"},
			{Address: "10.1.240.152", Hostname: "k8s-02", Description: "k8s node #2"},
			{Address: "10.1.240.153", Hostname: "k8s-03", Description: "k8s node #3"},
		}
		res := HostsResponse{
			Count: len(hosts),
			Hosts: hosts,
		}
		replyJSON(w, http.StatusOK, res)
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
