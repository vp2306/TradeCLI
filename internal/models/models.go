package models

import "time"

// Quote represents a stock quote
type Quote struct {
	Symbol    string
	BidPrice  float64
	AskPrice  float64
	LastPrice float64
	Volume    int64
	Timestamp time.Time
}

// Order represents a trade order
type Order struct {
	ID            string
	Symbol        string
	Quantity      int
	Side          string // "buy" or "sell"
	Type          string // "market" or "limit"
	Price         float64
	Status        string
	FilledQty     int
	FilledAvgPx   float64
	SubmittedAt   time.Time
	FilledAt      *time.Time
}

// Position represents a stock position
type Position struct {
	Symbol           string
	Quantity         int
	AvgEntryPrice    float64
	CurrentPrice     float64
	MarketValue      float64
	CostBasis        float64
	UnrealizedPL     float64
	UnrealizedPLPct  float64
	Side             string
}

// Account represents account information
type Account struct {
	ID                  string
	Cash                float64
	PortfolioValue      float64
	BuyingPower         float64
	Equity              float64
	LastEquity          float64
	DaytradeCount       int
	PatternDayTrader    bool
}