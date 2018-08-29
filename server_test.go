package ipmonitor

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	// for gorm to use sqlite
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type TestHostsValue struct {
	Name   string
	Method string
	Input  TestHostRequest
	Expect TestHostResponse
}

type TestHostValue struct {
	Name     string
	Method   string
	Endpoint string
	Expect   interface{}
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
	const dbfile = "unittest.db"
	err := OpenDB(dbfile)
	if err != nil {
		log.Fatalln(err)
	}
	InitDB()
	defer Conn.DB.Close()

	res := m.Run()
	err = os.Remove(dbfile)
	if err != nil {
		log.Println(err)
	}
	os.Exit(res)
}

func TestHostsHandler(t *testing.T) {
	router := NewHTTPHandler()

	values := []TestHostsValue{
		TestHostsValue{
			Name:   "Empty GET",
			Method: "GET",
			Input:  TestHostRequest{},
			Expect: TestHostResponse{http.StatusOK, HostsResponse{Count: 0, Hosts: []Host{}}},
		},
		TestHostsValue{
			Name:   "POST",
			Method: "POST",
			Input:  TestHostRequest{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"},
			Expect: TestHostResponse{http.StatusCreated, Host{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"}},
		},
		TestHostsValue{
			Name:   "GET 1 record",
			Method: "GET",
			Input:  TestHostRequest{},
			Expect: TestHostResponse{http.StatusOK, HostsResponse{Count: 1, Hosts: []Host{Host{Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"}}}},
		},
		TestHostsValue{
			Name:   "Invalid POST request body #1",
			Method: "POST",
			Input:  TestHostRequest{Address: "10.1.240.151"},
			Expect: TestHostResponse{http.StatusBadRequest, ErrorResponse{Status: http.StatusBadRequest, Message: "Key \"address\" and \"hostname\" are required"}},
		},
		TestHostsValue{
			Name:   "Invalid POST request body #2",
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
			t.Errorf("'%s': Wrong status code was returned. '%d' != '%d'", v.Name, res.StatusCode, v.Expect.StatusCode)
		}

		switch expected := v.Expect.Response.(type) {
		case HostsResponse:
			var hosts HostsResponse
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &hosts)

			if hosts.Count != expected.Count {
				t.Errorf("'%s': Wrong hosts count was returned. '%d' != '%d'", v.Name, hosts.Count, expected.Count)
			}
			for i := range hosts.Hosts {
				checkHostResponse(t, v.Name, hosts.Hosts[i], expected.Hosts[i])
			}
		case Host:
			var host Host
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &host)
			checkHostResponse(t, v.Name, host, expected)
		case ErrorResponse:
			var err ErrorResponse
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &err)
			checkErrorResponse(t, v.Name, err, expected)
		default:
		}
	}
}

func TestHostHandler(t *testing.T) {
	router := NewHTTPHandler()
	// init test data
	values := []TestHostValue{
		TestHostValue{
			Name:     "Non-Integer Key",
			Method:   "GET",
			Endpoint: "/hosts/x",
			Expect:   ErrorResponse{Status: http.StatusBadRequest, Message: "Key \"id\" must be integer"},
		},
		TestHostValue{
			Name:     "Not Found Key",
			Method:   "GET",
			Endpoint: "/hosts/9999",
			Expect:   ErrorResponse{Status: http.StatusNotFound, Message: "ID \"9999\" not found"},
		},
		TestHostValue{
			Name:     "Get first record",
			Method:   "GET",
			Endpoint: "/hosts/1",
			Expect:   Host{ID: 1, Address: "10.1.240.151", Hostname: "k8s-01", Description: "k8s node #1"},
		},
	}

	for _, v := range values {
		req := httptest.NewRequest(v.Method, v.Endpoint, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)
		res := rr.Result()

		switch expect := v.Expect.(type) {
		case ErrorResponse:
			var err ErrorResponse
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &err)
			checkErrorResponse(t, v.Name, err, expect)
		case Host:
			var host Host
			body, _ := ioutil.ReadAll(res.Body)
			json.Unmarshal(body, &host)
			checkHostResponse(t, v.Name, host, expect)
		}
	}
}

func checkHostResponse(t *testing.T, testcase string, res, expected Host) {
	if res.Address != expected.Address {
		t.Errorf("'%s': Wrong Address was returned. '%s' != '%s'", testcase, res.Address, expected.Address)
	}
	if res.Hostname != expected.Hostname {
		t.Errorf("'%s': Wrong Hostname was returned. '%s' != '%s'", testcase, res.Hostname, expected.Hostname)
	}
	if res.Description != expected.Description {
		t.Errorf("'%s': Wrong Description was returned. '%s' != '%s'", testcase, res.Description, expected.Description)
	}
}

func checkErrorResponse(t *testing.T, testcase string, res, expected ErrorResponse) {
	if res.Status != expected.Status {
		t.Errorf("'%s': Wrong status code was returned. '%d' != '%d'", testcase, res.Status, expected.Status)
	}
	if res.Message != expected.Message {
		t.Errorf("'%s': Wrong message was returned. '%s' != '%s'", testcase, res.Message, expected.Message)
	}
}
