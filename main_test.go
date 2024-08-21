package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gorilla/mux"
	"github.com/rainbowmga/timetravel/api"
	"github.com/rainbowmga/timetravel/service"
)

var testServer *httptest.Server

func TestMain(m *testing.M) {
	// Set up the test server
	sqliteService, err := service.NewSQLiteRecordService("./test_records.db")
	if err != nil {
		fmt.Printf("Failed to create SQLite service: %v\n", err)
		os.Exit(1)
	}

	apiHandler := api.NewAPI(sqliteService)
	router := mux.NewRouter()
	apiRoute := router.PathPrefix("/api/v1").Subrouter()
	apiHandler.CreateRoutes(apiRoute)

	testServer = httptest.NewServer(router)
	defer testServer.Close()

	// Run the tests
	code := m.Run()

	// Clean up
	os.Remove("./test_records.db")
	os.Exit(code)
}

func TestCreateRecord(t *testing.T) {
	payload := map[string]string{"name": "John Doe", "email": "john@example.com"}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(testServer.URL+"/api/v1/records/1", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["id"] != float64(1) {
		t.Errorf("Expected id 1; got %v", result["id"])
	}
}

func TestGetRecord(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/records/1")
	if err != nil {
		t.Fatalf("Failed to get record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["id"] != float64(1) {
		t.Errorf("Expected id 1; got %v", result["id"])
	}
}

func TestUpdateRecord(t *testing.T) {
	payload := map[string]string{"email": "johndoe@example.com"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/v1/records/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to update record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})
	if data["email"] != "johndoe@example.com" {
		t.Errorf("Expected updated email johndoe@example.com; got %v", data["email"])
	}
}

func TestDeleteField(t *testing.T) {
	payload := map[string]*string{"name": nil}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/v1/records/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("Failed to delete field: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	data := result["data"].(map[string]interface{})
	if _, exists := data["name"]; exists {
		t.Errorf("Expected name field to be deleted, but it still exists")
	}
}

func TestGetNonExistentRecord(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/records/999")
	if err != nil {
		t.Fatalf("Failed to get non-existent record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", resp.Status)
	}
}
