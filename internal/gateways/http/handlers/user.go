package handlers

import (
	"fmt"
	"homework/internal/domain"
	"homework/internal/gateways/http/models"
	"homework/internal/gateways/http/types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

func CreateUserHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.UserToCreate

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

		user := &domain.User{
			ID:   1,
			Name: *req.Name,
		}

		if _, err := u.User.RegisterUser(c.Request.Context(), user); err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{
				Reason: swag.String(err.Error()),
			})
			return
		}

		response := models.User{
			ID:   &user.ID,
			Name: &user.Name,
		}
		c.JSON(http.StatusOK, response)
	}
}

func BindSensorHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIdStr := c.Param("user_id")
		userId, err := strconv.ParseInt(userIdStr, 10, 64)
		if err != nil || userId <= 0 {
			c.JSON(http.StatusBadRequest, models.Error{Reason: swag.String("incorrect user_id")})
			return
		}

		var req models.SensorToUserBinding

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

		if err := u.User.AttachSensorToUser(c.Request.Context(), userId, *req.SensorID); err != nil {
			c.JSON(http.StatusNotFound, models.Error{
				Reason: swag.String(err.Error()),
			})
			return
		}

		c.Status(http.StatusCreated)
	}
}

func sensorsByUser(c *gin.Context, u types.UseCases, withBody bool) {
	userIdStr := c.Param("user_id")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusUnprocessableEntity, models.Error{Reason: swag.String("incorrect user_id")})
		return
	}

	sensors, err := u.User.GetUserSensors(c.Request.Context(), userId)
	if err != nil {
		fmt.Println(">>>" + err.Error() + "<<<")
		c.JSON(http.StatusNotFound, models.Error{Reason: swag.String(err.Error())})
		return
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

func GetSensorsByUserHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensorsByUser(c, u, true)
	}
}

func HeadSensorsByUserHandler(u types.UseCases) gin.HandlerFunc {
	return func(c *gin.Context) {
		sensorsByUser(c, u, false)
	}
}
