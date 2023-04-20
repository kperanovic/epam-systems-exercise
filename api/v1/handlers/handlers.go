package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kperanovic/epam-systems/api/v1/types"
	"github.com/kperanovic/epam-systems/internal/kafka"
	"github.com/kperanovic/epam-systems/internal/storage"
	"go.uber.org/zap"
)

type RESTHandlers struct {
	log      *zap.Logger
	store    storage.Storage
	producer *kafka.Producer
}

func NewRESTHandlers(log *zap.Logger, store storage.Storage, producer *kafka.Producer) *RESTHandlers {
	return &RESTHandlers{
		log:      log,
		store:    store,
		producer: producer,
	}
}

// HandleGetCompany handles the GET endpoint "/v1/company/".
// It will validate the request and fetch the data from storage.
func (h *RESTHandlers) HandleGetCompany(c *gin.Context) {
	id := c.Param("id")

	h.log.Info("received getCompany request", zap.String("id", id))

	company, err := h.store.GetCompany(uuid.MustParse(id))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "unable to fetch company",
			})

		return

	}

	if company == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"message": "company not found",
			})

		return
	}

	c.JSON(http.StatusOK, company)
}

// HandleCreateCompany handles the POST endpoint "/v1/company/".
// It will validate the request body and save the data in storage.
// On successfull save it will Send a Message in "company.commands" Kafka topic.
func (h *RESTHandlers) HandleCreateCompany(c *gin.Context) {
	var company types.Company

	// Check if required bindings are satisfied
	if err := c.ShouldBindJSON(&company); err != nil {
		h.log.Error("error binding request body", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "invalid request. Please check the request body",
			})

		return
	}

	h.log.Info("received createCompany request", zap.Any("req", company))

	if err := h.store.SaveCompany(&company); err != nil {
		h.log.Error("error saving company", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error occured while processing request",
		})

		return
	}

	if err := h.producer.SendMessage(context.TODO(), "company.commands", "", "COMPANY_CREATED", company); err != nil {
		h.log.Error("error sending kafka message", zap.Error(err))

		// Intentionally not returning an error since failing to send a message into kafka topic has nothing to do
		// with the behaviour of the rest handler.
	}

	c.JSON(http.StatusOK, nil)
}

// HandlePatchCOmpany handles the PATCH endpoint "/v1/company/:id".
// It will validate the request body and will update the existing data in storage.
// On successfull update it will send a message in "company.commands" Kafka topic.
func (h *RESTHandlers) HandlePatchCompany(c *gin.Context) {
	var company types.Company

	id := c.Param("id")

	// Check if required bindings are satisfied
	if err := c.ShouldBindJSON(&company); err != nil {
		h.log.Error("error binding request body", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusBadRequest,
			gin.H{
				"error":   err.Error(),
				"message": "invalid request. Please check the request body",
			})

		return
	}

	h.log.Info("received patchCompany request", zap.String("id", id), zap.Any("company", company))

	if err := h.store.UpdateCompany(uuid.MustParse(id), &company); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error occured while processing request",
		})

		return
	}

	if err := h.producer.SendMessage(context.TODO(), "company.commands", "", "COMPANY_UPDATED", company); err != nil {
		h.log.Error("error sending kafka message", zap.Error(err))

		// Intentionally not returning an error since failing to send a message into kafka topic has nothing to do
		// with the behaviour of the rest handler.
	}

	c.JSON(http.StatusOK, nil)
}

// HandleDeleteCompany handles the DELETE endpoint "/v1/company/:id".
// It will validate the request body and will delete the existing data in storage.
// On successfull delete it will send a message in "company.commands" Kafka topic.
func (h *RESTHandlers) HandleDeleteCompany(c *gin.Context) {
	id := c.Param("id")

	h.log.Info("received deleteCompany request", zap.String("id", id))

	if err := h.store.DeleteCompany(uuid.MustParse(id)); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"message": "error occured while processing request",
		})

		return
	}

	if err := h.producer.SendMessage(context.TODO(), "company.commands", "", "COMPANY_DELETED", id); err != nil {
		h.log.Error("error sending kafka message", zap.Error(err))

		// Intentionally not returning an error since failing to send a message into kafka topic has nothing to do
		// with the behaviour of the rest handler.
	}

	c.JSON(http.StatusOK, nil)
}
