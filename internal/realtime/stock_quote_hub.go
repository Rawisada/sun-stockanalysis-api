package realtime

import (
	"sync"

	"github.com/gofiber/websocket/v2"

	"sun-stockanalysis-api/internal/models"
)

type StockQuoteNotifier interface {
	NotifyQuote(quote *models.StockQuote)
}

type StockQuoteHub struct {
	mu    sync.Mutex
	conns map[*websocket.Conn]struct{}
}

func NewStockQuoteHub() *StockQuoteHub {
	return &StockQuoteHub{
		conns: make(map[*websocket.Conn]struct{}),
	}
}

func (h *StockQuoteHub) Register(conn *websocket.Conn) {
	if h == nil || conn == nil {
		return
	}
	h.mu.Lock()
	h.conns[conn] = struct{}{}
	h.mu.Unlock()
}

func (h *StockQuoteHub) Unregister(conn *websocket.Conn) {
	if h == nil || conn == nil {
		return
	}
	h.mu.Lock()
	delete(h.conns, conn)
	h.mu.Unlock()
}

func (h *StockQuoteHub) NotifyQuote(quote *models.StockQuote) {
	if h == nil || quote == nil {
		return
	}
	h.mu.Lock()
	defer h.mu.Unlock()

	payload := struct {
		Quote *models.StockQuote `json:"quote"`
	}{
		Quote: quote,
	}

	for conn := range h.conns {
		if err := conn.WriteJSON(payload); err != nil {
			_ = conn.Close()
			delete(h.conns, conn)
		}
	}
}
