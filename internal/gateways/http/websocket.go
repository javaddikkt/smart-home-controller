package http

import (
	"errors"
	"homework/internal/gateways/http/types"
	"homework/internal/usecase"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/coder/websocket"
	"github.com/coder/websocket/wsjson"
)

type WebSocketHandler struct {
	useCases    types.UseCases
	mu          sync.Mutex
	connections map[*websocket.Conn]struct{}
}

func NewWebSocketHandler(useCases types.UseCases) *WebSocketHandler {
	return &WebSocketHandler{
		useCases:    useCases,
		connections: make(map[*websocket.Conn]struct{}),
	}
}

func (h *WebSocketHandler) HandleEvents(c *gin.Context) {
	sensorParam := c.Param("sensor_id")
	sensorID, err := strconv.ParseInt(sensorParam, 10, 64)
	if err != nil {
		c.AbortWithStatus(400)
		return
	}

	ctx := c.Request.Context()
	_, err = h.useCases.Sensor.GetSensorByID(ctx, sensorID)
	if err != nil {
		if errors.Is(err, usecase.ErrSensorNotFound) {
			c.AbortWithStatus(404)
			return
		}
		c.AbortWithStatus(500)
		return
	}

	conn, err := websocket.Accept(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	h.mu.Lock()
	h.connections[conn] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.connections, conn)
		h.mu.Unlock()
	}()

	h.Handle(c, conn, sensorID)
}

func (h *WebSocketHandler) Handle(c *gin.Context, conn *websocket.Conn, sensorID int64) {
	ctx := c.Request.Context()

	go func() {
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				return
			}
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			conn.Close(websocket.StatusNormalClosure, "")
			return
		case <-ticker.C:
			event, err := h.useCases.Event.GetLastEventBySensorID(ctx, sensorID)
			if err != nil {
				continue
			}
			if event == nil {
				continue
			}
			if err := wsjson.Write(ctx, conn, event); err != nil {
				conn.Close(websocket.StatusInternalError, "")
				return
			}
		}
	}
}

func (h *WebSocketHandler) Shutdown() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	for conn := range h.connections {
		conn.Close(websocket.StatusNormalClosure, "")
	}
	return nil
}
