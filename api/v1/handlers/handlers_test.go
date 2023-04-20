package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/Shopify/sarama/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/google/uuid"
	middleware "github.com/kperanovic/epam-systems/api/v1/auth"
	"github.com/kperanovic/epam-systems/api/v1/types"
	"github.com/kperanovic/epam-systems/internal/kafka"
	"github.com/kperanovic/epam-systems/internal/logger"
	"github.com/kperanovic/epam-systems/internal/storage"
	"github.com/kperanovic/epam-systems/internal/token"
	"go.uber.org/zap"
)

func GinRouter() *gin.Engine {
	router := gin.Default()
	return router
}

func generateCompany() *types.Company {
	uid, _ := uuid.NewRandom()

	return &types.Company{
		ID:          uid,
		Name:        "test-company",
		Description: "description",
		Employees:   10,
		Registered:  true,
		CompanyType: 1,
	}
}

func TestNewRESTHandlers(t *testing.T) {
	// Create new dev logger
	log := logger.NewDevelopment()

	// Set new kafka mock producer
	producer := kafka.NewMockProducer(mocks.NewSyncProducer(t, nil), log)

	type args struct {
		log      *zap.Logger
		store    storage.Storage
		producer *kafka.Producer
	}
	tests := []struct {
		name string
		args args
		want *RESTHandlers
	}{
		{
			name: "Test NewRESTHandlers()",
			args: args{
				log:      log,
				store:    storage.NewMemoryStorage(),
				producer: producer,
			},
			want: &RESTHandlers{
				log:      log,
				store:    storage.NewMemoryStorage(),
				producer: producer,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRESTHandlers(tt.args.log, tt.args.store, tt.args.producer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRESTHandlers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRESTHandlers_HandleGetCompany(t *testing.T) {
	// Define new gin router
	r := GinRouter()

	// Set new dev logger
	log := logger.NewDevelopment()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mocks.NewSyncProducer(t, nil), log),
	)

	// generate new *types.Company struct
	company := generateCompany()

	// Save it in storage
	err := h.store.SaveCompany(company)
	assert.Equal(t, err, nil)

	// Define route
	r.GET("/v1/company/:id", h.HandleGetCompany)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/v1/company/%s", company.ID), nil)

	// Parse company to json
	jsonParse, _ := json.Marshal(company)

	// send request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Body.String() != string(jsonParse) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			w.Body.String(), string(jsonParse))
	}
	assert.Equal(t, http.StatusOK, w.Code)

	// Test if invalid id is sent
	req, _ = http.NewRequest("GET", fmt.Sprintf("/v1/company/%s", uuid.New()), nil)

	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRESTHandlers_HandleCreateCompany(t *testing.T) {
	// Define new gin router
	r := GinRouter()

	// Define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// generate new *types.Company struct
	company := generateCompany()

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Marshal into json
	jsonValue, _ := json.Marshal(company)

	// define a route
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.POST("/", h.HandleCreateCompany)

	req, _ := http.NewRequest("POST", "/v1/company/", bytes.NewBuffer(jsonValue))

	// create a token
	token, err := j.CreateToken(uuid.New(), company.Name, 10*time.Second)
	assert.Equal(t, err, nil)

	// Add token to Bearer header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// make request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRESTHandlers_HandleCreateCompany_BadRequest(t *testing.T) {
	// Define new gin router
	r := GinRouter()

	// Define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// create new invalid company struct
	company := &types.Company{
		ID:   uuid.New(),
		Name: "a random name longer than 15 characters",
	}

	// Create token
	token, err := j.CreateToken(uuid.New(), company.Name, 10*time.Second)
	assert.Equal(t, err, nil)

	// Generate an empty company struct and send it on the endpoint.
	jsonValue, _ := json.Marshal(company)

	// Define the route
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.POST("/", h.HandleCreateCompany)

	req, _ := http.NewRequest("POST", "/v1/company/", bytes.NewBuffer(jsonValue))

	// Add token to Bearer Header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRESTHandlers_HandleCreateCompanyWithoutToken(t *testing.T) {
	// define new gin router
	r := GinRouter()

	// define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// generate new *types.Company struct
	company := generateCompany()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// Generate the token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Define the endpoint
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.POST("/", h.HandleCreateCompany)

	// Marshal the request body
	jsonValue, _ := json.Marshal(company)

	// Define the request
	req, _ := http.NewRequest("POST", "/v1/company/", bytes.NewBuffer(jsonValue))

	// Make request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRESTHandlers_HandlePatchCompany(t *testing.T) {
	// define the gin router
	r := GinRouter()

	// define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// generate new *types.Company struct
	company := generateCompany()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// First we save the original value in storage
	h.store.SaveCompany(company)

	// make new company struct with different data
	patched := generateCompany()
	patched.Employees = 500
	patched.Description = "This is a changed description"

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Marshal the company value to json
	jsonValue, _ := json.Marshal(patched)

	// Declare PATCH endpoint and generate a request
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.PATCH("/:id", h.HandlePatchCompany).Use(middleware.AuthMiddleware(j))
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/company/%s", company.ID), bytes.NewBuffer(jsonValue))

	// Create JWT auth token
	token, err := j.CreateToken(uuid.New(), company.Name, 10*time.Second)
	assert.Equal(t, err, nil)

	// Add token to Bearer header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make a request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Get patched data from storage
	got, _ := h.store.GetCompany(company.ID)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEqual(t, company.Employees, got.Employees)
	assert.NotEqual(t, company.Description, got.Description)
}

func TestRESTHandlers_HandlePatchCompany_WithoutToken(t *testing.T) {
	// define new gin router
	r := GinRouter()

	// define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// generate new *types.Company struct
	company := generateCompany()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Marshal the company value to json
	jsonValue, _ := json.Marshal(company)

	// Declare PATCH endpoint and generate a request
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.PATCH("/:id", h.HandlePatchCompany).Use(middleware.AuthMiddleware(j))
	req, _ := http.NewRequest("PATCH", fmt.Sprintf("/v1/company/%s", company.ID), bytes.NewBuffer(jsonValue))

	// Make a request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRESTHandlers_HandleDeleteCompany(t *testing.T) {
	// define new gin router
	r := GinRouter()

	// define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// generate new *types.Company struct
	company := generateCompany()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// save company in storage
	h.store.SaveCompany(company)

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Marshal the company value to json
	jsonValue, _ := json.Marshal(company)

	// Declare PATCH endpoint and generate a request
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.DELETE("/:id", h.HandleDeleteCompany).Use(middleware.AuthMiddleware(j))
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/company/%s", company.ID), bytes.NewBuffer(jsonValue))

	// Create JWT auth token
	token, err := j.CreateToken(uuid.New(), company.Name, 10*time.Second)
	assert.Equal(t, err, nil)

	// Add token to Bearer header
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	// Make a request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Try to fetch deleted company from store
	got, err := h.store.GetCompany(company.ID)

	assert.Equal(t, err, nil)
	assert.Equal(t, got, nil)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestRESTHandlers_HandleDeleteCompany_WithoutToken(t *testing.T) {
	// define new gin router
	r := GinRouter()

	// define new dev logger
	log := logger.NewDevelopment()

	// Define new kafka mock producer and set mock expectation
	// so that mock producer knows how to behave when it needs to send message.
	mockProducer := mocks.NewSyncProducer(t, nil)
	mockProducer.ExpectSendMessageAndSucceed()

	// generate new *types.Company struct
	company := generateCompany()

	// Initiate new RESTHandlers struct
	h := NewRESTHandlers(
		log,
		storage.NewMemoryStorage(),
		kafka.NewMockProducer(mockProducer, log),
	)

	// save company in storage
	h.store.SaveCompany(company)

	// Generate a token
	j, err := token.NewJWTToken("KLguRWx03zXcWwDXywrxgwTS7r39QaF1")
	assert.Equal(t, err, nil)

	// Marshal the company value to json
	jsonValue, _ := json.Marshal(company)

	// Declare PATCH endpoint and generate a request
	g := r.Group("/v1/company").Use(middleware.AuthMiddleware(j))
	g.DELETE("/:id", h.HandleDeleteCompany).Use(middleware.AuthMiddleware(j))
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/v1/company/%s", company.ID), bytes.NewBuffer(jsonValue))

	// Make a request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
