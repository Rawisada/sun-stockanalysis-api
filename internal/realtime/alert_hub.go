package realtime

import (
	"sync"

	"github.com/gofiber/websocket/v2"

	"sun-stockanalysis-api/internal/models"
)

type AlertEventNotifier interface {
	Notify(event *models.AlertEvent, message string)
}

type AlertHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
}

func NewAlertHub() *AlertHub {
	return &AlertHub{
		conns: make(map[*websocket.Conn]struct{}),
	}
}

func (h *AlertHub) Register(conn *websocket.Conn) {
	if h == nil || conn == nil {
		return
	}
	h.mu.Lock()
	h.conns[conn] = struct{}{}
	h.mu.Unlock()
}

func (h *AlertHub) Unregister(conn *websocket.Conn) {
	if h == nil || conn == nil {
		return
	}
	h.mu.Lock()
	delete(h.conns, conn)
	h.mu.Unlock()
}

func (h *AlertHub) Notify(event *models.AlertEvent, message string) {
	if h == nil || event == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	payload := struct {
		Event   *models.AlertEvent `json:"event"`
		Message string            `json:"message"`
	}{
		Event:   event,
		Message: message,
	}

	for conn := range h.conns {
		if err := conn.WriteJSON(payload); err != nil {
			_ = conn.Close()
			delete(h.conns, conn)
		}
	}
}
