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
	apiHandler.CreateRoutes(router)

	testServer = httptest.NewServer(router)
	defer testServer.Close()

	// Run the tests
	code := m.Run()

	// Clean up
	os.Remove("./test_records.db")
	os.Exit(code)
}

// V1 API Tests

func TestCreateRecordV1(t *testing.T) {
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

func TestGetRecordV1(t *testing.T) {
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

func TestUpdateRecordV1(t *testing.T) {
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

func TestDeleteFieldV1(t *testing.T) {
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

func TestGetNonExistentRecordV1(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v1/records/999")
	if err != nil {
		t.Fatalf("Failed to get non-existent record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status Bad Request; got %v", resp.Status)
	}
}

// V2 API Tests

func TestCreateRecordV2(t *testing.T) {
	payload := map[string]string{"name": "Jane Doe", "email": "jane@example.com"}
	body, _ := json.Marshal(payload)
	resp, err := http.Post(testServer.URL+"/api/v2/records/2", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["id"] != float64(2) {
		t.Errorf("Expected id 2; got %v", result["id"])
	}
	if result["version"] != float64(1) {
		t.Errorf("Expected version 1; got %v", result["version"])
	}
}

func TestGetRecordV2(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v2/records/2")
	if err != nil {
		t.Fatalf("Failed to get record: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["id"] != float64(2) {
		t.Errorf("Expected id 2; got %v", result["id"])
	}
	if result["version"] != float64(1) {
		t.Errorf("Expected version 1; got %v", result["version"])
	}
}

func TestUpdateRecordV2(t *testing.T) {
	payload := map[string]string{"email": "janedoe@example.com"}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", testServer.URL+"/api/v2/records/2", bytes.NewBuffer(body))
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
	if result["version"] != float64(2) {
		t.Errorf("Expected version 2; got %v", result["version"])
	}
	data := result["data"].(map[string]interface{})
	if data["email"] != "janedoe@example.com" {
		t.Errorf("Expected updated email janedoe@example.com; got %v", data["email"])
	}
}

func TestGetRecordVersionV2(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v2/records/2?version=1")
	if err != nil {
		t.Fatalf("Failed to get record version: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if result["id"] != float64(2) {
		t.Errorf("Expected id 2; got %v", result["id"])
	}
	if result["version"] != float64(1) {
		t.Errorf("Expected version 1; got %v", result["version"])
	}
}

func TestGetRecordVersionsV2(t *testing.T) {
	resp, err := http.Get(testServer.URL + "/api/v2/records/2/versions")
	if err != nil {
		t.Fatalf("Failed to get record versions: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK; got %v", resp.Status)
	}

	var result []int
	json.NewDecoder(resp.Body).Decode(&result)
	if len(result) != 2 {
		t.Errorf("Expected 2 versions; got %v", len(result))
	}
	if result[0] != 1 || result[1] != 2 {
		t.Errorf("Expected versions [1, 2]; got %v", result)
	}
}
