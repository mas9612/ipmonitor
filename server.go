package ipmonitor

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// NewHTTPHandler returns *mux.Router
func NewHTTPHandler() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hosts", hostsHandler).Methods("GET", "POST")
	r.HandleFunc("/hosts/{id}", hostHandler).Methods("GET", "PUT", "DELETE")

	return r
}

func hostsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		var hosts []Host
		err := Conn.DB.Find(&hosts).Error
		if err != nil {
			log.Println("[ERROR] /hosts: failed to get host records:", err)
			return
		}

		replyJSON(w, http.StatusOK, HostsResponse{
			Count: len(hosts),
			Hosts: hosts,
		})
		return
	case http.MethodPost:
		var host Host
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&host)
		if err != nil {
			replyError(w, http.StatusBadRequest, "Received malicious request")
			return
		}
		if host.Address == "" || host.Hostname == "" {
			replyError(w, http.StatusBadRequest, "Key \"address\" and \"hostname\" are required")
			return
		}

		result := Conn.DB.Create(&host)
		if result.Error != nil {
			replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
			return
		}
		replyJSON(w, http.StatusCreated, result.Value)
		return
	}
	replyError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s is not allowed in this URL", r.Method))
}

func hostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		replyError(w, http.StatusBadRequest, "Key \"id\" must be integer")
		return
	}

	switch r.Method {
	case http.MethodGet:
		var host Host
		err = Conn.DB.Where("id = ?", id).Find(&host).Error
		if err != nil {
			replyError(w, http.StatusNotFound, fmt.Sprintf("ID \"%d\" not found", id))
			return
		}
		replyJSON(w, http.StatusOK, Host{ID: host.ID, Address: host.Address, Hostname: host.Hostname, Description: host.Description})
		return
	case http.MethodPut:
		var host Host
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&host)
		if err != nil {
			replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
			return
		}
		if host.Address == "" || host.Hostname == "" {
			replyError(w, http.StatusBadRequest, "Key \"address\" and \"hostname\" are required")
			return
		}

		vars := mux.Vars(r)
		var now Host
		err = Conn.DB.Where("id = ?", vars["id"]).Find(&now).Error
		if err != nil && gorm.IsRecordNotFoundError(err) {
			replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
			return
		}
		now.Address = host.Address
		now.Hostname = host.Hostname
		now.Description = host.Description
		Conn.DB.Save(&now)

		replyJSON(w, http.StatusOK, now)
		return
	case http.MethodDelete:
		vars := mux.Vars(r)
		var host Host
		err = Conn.DB.Where("id = ?", vars["id"]).Find(&host).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
				return
			}
			replyError(w, http.StatusNotFound, fmt.Sprintf("ID \"%s\" not found", vars["id"]))
			return
		}
		Conn.DB.Delete(&host)
		replyJSON(w, http.StatusNoContent, nil)
		return
	}
	replyError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Method %s is not allowed in this URL", r.Method))
}

func replyJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if data != nil {
		encoder := json.NewEncoder(w)
		err := encoder.Encode(data)
		if err != nil {
			replyError(w, http.StatusInternalServerError, "Internal Server Error occured.")
		}
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
