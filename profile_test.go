package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/ericzty/eve/internal/db"
	"github.com/google/uuid"
)

const (
	host                  = "http://localhost:3000"
	dbUrl                 = "postgres://postgres:test@localhost:5432/postgres"
	testAdminName         = "Admin Test"
	testAdminEmail        = "admin@testing.com"
	testAdminPassword     = "adminPasswordTest"
	testAdminPasswordHash = "$2a$11$MCwHsWkdVATPJ0URdUFg9uvY6UdskKO.Mwc3Y2e9LKi.5GQFOhTCq"

	testUserName     = "User Test"
	testUserEmail    = "eric@testing.com"
	testUserPassword = "userPasswordTest"
)

var (
	testAdminId = uuid.New()
)

func TestSetup(t *testing.T) {
	// Insert test user and hypervisor
	ctx := context.Background()
	db.Init(dbUrl)
	_, err := db.Pool.Exec(
		ctx,
		"INSERT INTO profile (id, name, email, password, is_admin) VALUES ($1, $2, $3, $4, $5)",
		testAdminId, testAdminName, testAdminEmail, testAdminPasswordHash, true,
	)
	if err != nil {
		t.Fatal(err)
	}
}

func TestHealthCheck(t *testing.T) {
	resp, err := http.Get(host + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}
}

var adminToken string

func TestAdminLogin(t *testing.T) {
	request := map[string]string{
		"email":    testAdminEmail,
		"password": testAdminPassword,
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(host+"/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	adminToken = response["token"]
	if adminToken == "" {
		t.Fatal("token is empty")
	}
}

func TestAdminGetProfile(t *testing.T) {
	req, err := http.NewRequest("GET", host+"/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response["id"] != testAdminId.String() {
		t.Fatalf("expected %v; got %v", testAdminId, response["id"])
	}
	if response["name"] != testAdminName {
		t.Fatalf("expected %v; got %v", testAdminName, response["name"])
	}
	if response["email"] != testAdminEmail {
		t.Fatalf("expected %v; got %v", testAdminEmail, response["email"])
	}
	if (response["lastLogin"] == nil) || (response["created"] == nil) || (response["updated"] == nil) {
		t.Fatal("lastLogin, created, and/or updated is nil")
	}
}

var testUserId string

func TestCreateUser(t *testing.T) {
	request := map[string]string{
		"name":     testUserName,
		"email":    testUserEmail,
		"password": testUserPassword,
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest("POST", host+"/admin/users", bytes.NewBuffer(reqBody))
	req.Header.Set("Authorization", "Bearer "+adminToken)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected Created; got %v, body: %v", resp.Status, respBody.String())
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	testUserId = response["uuid"]
	if testUserId == "" {
		t.Fatal("uuid is empty")
	}
}

var userToken string

func TestUserLogin(t *testing.T) {
	request := map[string]string{
		"email":    testUserEmail,
		"password": testUserPassword,
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(host+"/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}

	var response map[string]string
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	userToken = response["token"]
	if userToken == "" {
		t.Fatal("token is empty")
	}
}

func TestUserGetProfile(t *testing.T) {
	req, err := http.NewRequest("GET", host+"/users/me", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+userToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response["id"] != testUserId {
		t.Fatalf("expected %v; got %v", testUserId, response["id"])
	}
	if response["name"] != testUserName {
		t.Fatalf("expected %v; got %v", testUserName, response["name"])
	}
	if response["email"] != testUserEmail {
		t.Fatalf("expected %v; got %v", testUserEmail, response["email"])
	}
	if (response["lastLogin"] == nil) || (response["created"] == nil) || (response["updated"] == nil) {
		t.Fatal("lastLogin, created, and/or updated is nil")
	}
}

func TestAdminLogout(t *testing.T) {
	req, err := http.NewRequest("POST", host+"/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+adminToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response["message"] != "logout success" {
		t.Fatalf("expected %v; got %v", "logout success", response["message"])
	}
}

func TestUserLogout(t *testing.T) {
	req, err := http.NewRequest("POST", host+"/logout", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+userToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		respBody := new(bytes.Buffer)
		respBody.ReadFrom(resp.Body)
		t.Fatalf("expected OK; got %v, body: %v", resp.Status, respBody.String())
	}
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		t.Fatal(err)
	}
	if response["message"] != "logout success" {
		t.Fatalf("expected %v; got %v", "logout success", response["message"])
	}
}

func TestCleanUp(t *testing.T) {
	//TODO: Prlly remove this later
	ctx := context.Background()
	_, err := db.Pool.Exec(ctx, "DELETE FROM sessions WHERE owner = $1", testAdminId)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Pool.Exec(ctx, "DELETE FROM profile WHERE id = $1", testAdminId)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Pool.Exec(ctx, "DELETE FROM sessions WHERE owner = $1", testUserId)
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Pool.Exec(ctx, "DELETE FROM profile WHERE id = $1", testUserId)
	if err != nil {
		t.Fatal(err)
	}
	db.Pool.Close()
}
