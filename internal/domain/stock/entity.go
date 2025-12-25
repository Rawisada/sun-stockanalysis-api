package stock

import "time"

type Stock struct {
    ID        string
    Symbol    string
	Company   string
    Price     float64
    CreatedAt time.Time
    UpdatedAt time.Time
}
