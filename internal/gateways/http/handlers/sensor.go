package handlers

import (
	"homework/internal/domain"
	"homework/internal/gateways/http/models"
	"homework/internal/gateways/http/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func CreateSensorHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.SensorToCreate

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

		sensor := &domain.Sensor{
			ID:           1,
			SerialNumber: *req.SerialNumber,
			Type:         domain.SensorType(*req.Type),
			CurrentState: 0,
			Description:  *req.Description,
			IsActive:     *req.IsActive,
			RegisteredAt: time.Now(),
			LastActivity: time.Now(),
		}

		if _, err := u.Sensor.RegisterSensor(c.Request.Context(), sensor); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Reason: swag.String(err.Error()),
			})
			return
		}

		c.JSON(http.StatusOK, mapDomainSensorToModel(sensor))
	}
}

func GetSensorsHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensors(c, u, true)
	}
}

func HeadSensorsHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensors(c, u, false)
	}
}

func GetSensorByIdHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensorByID(c, u, true)
	}
}

func HeadSensorByIdHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensorByID(c, u, false)
	}
}

func GetSensorHistoryHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensorIDStr := c.Param("sensor_id")
		sensorID, err := strconv.ParseInt(sensorIDStr, 10, 64)
		if err != nil || sensorID <= 0 {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Reason: swag.String("invalid sensor_id"),
			})
			return
		}

		startDateStr := c.Query("start_date")
		startDate, err := time.Parse(time.RFC3339, startDateStr)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Reason: swag.String("invalid start_date"),
			})
		}

		endDateStr := c.Query("end_date")
		endDate, err := time.Parse(time.RFC3339, endDateStr)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, models.Error{
				Reason: swag.String("invalid end_date"),
			})
		}

		events, err := u.Event.GetEventsInRangeBySensorID(c.Request.Context(), sensorID, startDate, endDate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Reason: swag.String(err.Error()),
			})
		}

		resp := make([]*models.SensorHistoryItem, 0, len(events))
		for _, e := range events {
			item := &models.SensorHistoryItem{
				Payload:   e.Payload,
				Timestamp: strfmt.DateTime(e.Timestamp),
			}
			resp = append(resp, item)
		}

		c.JSON(http.StatusOK, resp)
	}
}

func sensors(c *gin.Context, u types.UseCases, withBody bool) {
	sensors, err := u.Sensor.GetSensors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{
			Reason: swag.String(err.Error()),
		})
	}
	resp := make([]models.Sensor, 0, len(sensors))
	for _, s := range sensors {
		sensor := mapDomainSensorToModel(&s)
		resp = append(resp, sensor)
	}

	if withBody {
		c.JSON(http.StatusOK, resp)
	} else {
		c.Header("Content-Length", "0")
		c.Status(http.StatusOK)
	}
}

func sensorByID(c *gin.Context, u types.UseCases, withBody bool) {
	sensorIDStr := c.Param("sensor_id")
	sensorID, err := strconv.ParseInt(sensorIDStr, 10, 64)
	if err != nil || sensorID <= 0 {
		c.JSON(http.StatusUnprocessableEntity, models.Error{
			Reason: swag.String("invalid sensor_id"),
		})
		return
	}

	sensor, err := u.Sensor.GetSensorByID(c.Request.Context(), sensorID)
	if err != nil {
		c.JSON(http.StatusNotFound, models.Error{
			Reason: swag.String("sensor not found"),
		})
		return
	}

	if withBody {
		c.JSON(http.StatusOK, mapDomainSensorToModel(sensor))
	} else {
		c.Header("Content-Length", "0")
		c.Status(http.StatusOK)
	}
}

func mapDomainSensorToModel(s *domain.Sensor) models.Sensor {
	return models.Sensor{
		ID:           &s.ID,
		SerialNumber: &s.SerialNumber,
		Type:         swag.String(string(s.Type)),
		Description:  &s.Description,
		IsActive:     &s.IsActive,
		CurrentState: &s.CurrentState,
		RegisteredAt: (*strfmt.DateTime)(&s.RegisteredAt),
		LastActivity: (*strfmt.DateTime)(&s.LastActivity),
	}
}
