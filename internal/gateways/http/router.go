package http

import (
	"homework/internal/gateways/http/handlers"
	"homework/internal/gateways/http/types"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func setupRouter(r *gin.Engine, u types.UseCases, ws *WebSocketHandler) {
	r.HandleMethodNotAllowed = true
	r.NoMethod(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusMethodNotAllowed)
	})

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/events", RequireJSONContentType(), handlers.CreateEventHandler(u))
	r.OPTIONS("/events", optionsHandler([]string{"POST", "OPTIONS"}))

	r.GET("/sensors", RequireJSONAccept(), handlers.GetSensorsHandler(u))
	r.HEAD("/sensors", RequireJSONAccept(), handlers.HeadSensorsHandler(u))
	r.POST("/sensors", RequireJSONContentType(), handlers.CreateSensorHandler(u))
	r.OPTIONS("/sensors", optionsHandler([]string{"GET", "HEAD", "POST", "OPTIONS"}))

	r.GET("/sensors/:sensor_id/events", ws.HandleEvents)

	r.GET("/sensors/:sensor_id", RequireJSONAccept(), handlers.GetSensorByIdHandler(u))
	r.HEAD("/sensors/:sensor_id", RequireJSONAccept(), handlers.HeadSensorByIdHandler(u))
	r.OPTIONS("/sensors/:sensor_id", optionsHandler([]string{"GET", "HEAD", "OPTIONS"}))

	r.GET("/sensors/:sensor_id/history", RequireJSONAccept(), handlers.GetSensorHistoryHandler(u))

	r.POST("/users", RequireJSONContentType(), handlers.CreateUserHandler(u))
	r.OPTIONS("/users", optionsHandler([]string{"POST", "OPTIONS"}))

	r.GET("/users/:user_id/sensors", RequireJSONAccept(), handlers.GetSensorsByUserHandler(u))
	r.HEAD("/users/:user_id/sensors", RequireJSONAccept(), handlers.HeadSensorsByUserHandler(u))
	r.POST("/users/:user_id/sensors", RequireJSONContentType(), handlers.BindSensorHandler(u))
	r.OPTIONS("/users/:user_id/sensors", optionsHandler([]string{"GET", "HEAD", "POST", "OPTIONS"}))
}

func optionsHandler(methods []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Allow", strings.Join(methods, ","))
		c.Status(http.StatusNoContent)
	}
}
