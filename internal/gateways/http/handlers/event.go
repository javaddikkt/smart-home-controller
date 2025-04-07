package handlers

import (
	"homework/internal/domain"
	"homework/internal/gateways/http/models"
	"homework/internal/gateways/http/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func CreateEventHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SensorEvent

		if err := c.ShouldBindJSON(&req); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		if err := req.Validate(strfmt.Default); err != nil {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Reason: swag.String(err.Error()),
			})
			return
		}

		event := &domain.Event{
			Timestamp:          time.Now(),
			SensorSerialNumber: *req.SensorSerialNumber,
			SensorID:           0,
			Payload:            *req.Payload,
		}

		if err := u.Event.ReceiveEvent(c.Request.Context(), event); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Reason: swag.String(err.Error()),
			})
			return
		}

		response := models.SensorEvent{
			Payload:            &event.Payload,
			SensorSerialNumber: &event.SensorSerialNumber,
		}
		c.JSON(http.StatusCreated, response)
	}
}
