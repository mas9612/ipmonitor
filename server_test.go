package ipmonitor

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	// for gorm to use sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type TestHostsValue struct {
	Method string
	Input  TestHostRequest
	Expect TestHostResponse
}

type TestHostRequest struct {
	Address     string `json:"address"`
	Hostname    string `json:"hostname"`
	Description string `json:"description"`
}

type TestHostResponse struct {
	StatusCode int
	Response   interface{}
}

func TestMain(m *testing.M) {
	const DBfile = "unittest.db"
	db, err := gorm.Open("sqlite3", DBfile)
	if err != nil {
		log.Fatalln("Failed to open DB")
	}
	defer db.Close()

	db.AutoMigrate(&Host{})

	res := m.Run()
	err = os.Remove(DBfile)
	if err != nil {
		log.Println(err)
	}
	os.Exit(res)
}

func TestHostsHandler(t *testing.T) {
	// TODO: fix to use "unittest.db"
	// TODO: add test which compare response
	router := NewHTTPHandler()

	values := []TestHostsValue{
		TestHostsValue{
			Method: "GET",
			Input:  TestHostRequest{},
			Expect: TestHostResponse{http.StatusOK, HostsResponse{Count: 0, Hosts: []Host{}}},
		},
		TestHostsValue{
			Method: "POST",
			Input:  TestHostRequest{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"},
			Expect: TestHostResponse{http.StatusCreated, HostsResponse{Count: 1, Hosts: []Host{Host{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"}}}},
		},
		TestHostsValue{
			Method: "GET",
			Input:  TestHostRequest{},
			Expect: TestHostResponse{http.StatusOK, HostsResponse{Count: 1, Hosts: []Host{Host{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"}}}},
		},
		TestHostsValue{
			Method: "POST",
			Input:  TestHostRequest{Address: "10.1.240.151"},
			Expect: TestHostResponse{http.StatusBadRequest, ErrorResponse{Status: http.StatusBadRequest, Message: "Key \"address\" and \"hostname\" are required"}},
		},
		TestHostsValue{
			Method: "POST",
			Input:  TestHostRequest{Hostname: "k8s-02"},
			Expect: TestHostResponse{http.StatusBadRequest, ErrorResponse{Status: http.StatusBadRequest, Message: "Key \"address\" and \"hostname\" are required"}},
		},
	}

	for _, v := range values {
		payload, _ := json.Marshal(v.Input)
		rr := httptest.NewRecorder()
		var req *http.Request
		if v.Method == "GET" {
			req = httptest.NewRequest(v.Method, "/hosts", nil)
		} else {
			req = httptest.NewRequest(v.Method, "/hosts", bytes.NewReader(payload))
		}
		router.ServeHTTP(rr, req)

		res := rr.Result()
		if res.StatusCode != v.Expect.StatusCode {
			t.Errorf("GET /hosts: Wrong status code was returned. '%d' != '%d'", res.StatusCode, v.Expect.StatusCode)
		}
	}
}

func TestHostHandler(t *testing.T) {

}

func checkHostsResponse(t *testing.T, res, expected Host) {
	if res.Address != expected.Address {
		t.Errorf("GET /hosts: Wrong Address was returned. '%s' != '%s'", res.Address, expected.Address)
	}
	if res.Hostname != expected.Hostname {
		t.Errorf("GET /hosts: Wrong Hostname was returned. '%s' != '%s'", res.Hostname, expected.Hostname)
	}
	if res.Description != expected.Description {
		t.Errorf("GET /hosts: Wrong Description was returned. '%s' != '%s'", res.Description, expected.Description)
	}
}
