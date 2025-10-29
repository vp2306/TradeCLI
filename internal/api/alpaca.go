package api

import (
	"fmt"

	"github.com/alpacahq/alpaca-trade-api-go/v3/alpaca"
	"github.com/alpacahq/alpaca-trade-api-go/v3/marketdata"
	"github.com/shopspring/decimal"
	"github.com/vp2306/terminal-trader/internal/config"
	"github.com/vp2306/terminal-trader/internal/models"
)

// Client wraps the Alpaca API client
type Client struct {
	alpaca     *alpaca.Client
	marketdata *marketdata.Client
	config     *config.Config
}

// NewClient creates a new Alpaca API client
func NewClient(cfg *config.Config) *Client {
	alpacaClient := alpaca.NewClient(alpaca.ClientOpts{
		APIKey:    cfg.Alpaca.APIKey,
		APISecret: cfg.Alpaca.APISecret,
		BaseURL:   cfg.Alpaca.BaseURL,
	})

	mdClient := marketdata.NewClient(marketdata.ClientOpts{
		APIKey:    cfg.Alpaca.APIKey,
		APISecret: cfg.Alpaca.APISecret,
	})

	return &Client{
		alpaca:     alpacaClient,
		marketdata: mdClient,
		config:     cfg,
	}
}

// GetAccount retrieves account information
func (c *Client) GetAccount() (*models.Account, error) {
	acct, err := c.alpaca.GetAccount()
	if err != nil {
		return nil, fmt.Errorf("failed to get account: %w", err)
	}

	return &models.Account{
		ID:                  acct.ID,
		Cash:                acct.Cash.InexactFloat64(),
		PortfolioValue:      acct.PortfolioValue.InexactFloat64(),
		BuyingPower:         acct.BuyingPower.InexactFloat64(),
		Equity:              acct.Equity.InexactFloat64(),
		LastEquity:          acct.LastEquity.InexactFloat64(),
		DaytradeCount:       int(acct.DaytradeCount),
		PatternDayTrader:    acct.PatternDayTrader,
	}, nil
}

// GetQuote retrieves a stock quote
func (c *Client) GetQuote(symbol string) (*models.Quote, error) {
	trade, err := c.marketdata.GetLatestTrade(symbol, marketdata.GetLatestTradeRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get quote for %s: %w", symbol, err)
	}

	return &models.Quote{
		Symbol:    symbol,
		LastPrice: trade.Price,
		Volume:    int64(trade.Size),
		Timestamp: trade.Timestamp,
	}, nil
}

// PlaceOrder places a market order
// PlaceOrder places a market order
func (c *Client) PlaceOrder(symbol string, qty int, side string) (*models.Order, error) {
	orderSide := alpaca.Buy
	if side == "sell" {
		orderSide = alpaca.Sell
	}

	qtyDecimal := decimal.NewFromInt(int64(qty))
	
	order, err := c.alpaca.PlaceOrder(alpaca.PlaceOrderRequest{
		Symbol:      symbol,
		Qty:         &qtyDecimal,
		Side:        orderSide,
		Type:        alpaca.Market,
		TimeInForce: alpaca.Day,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to place order: %w", err)
	}

	return &models.Order{
		ID:          order.ID,
		Symbol:      order.Symbol,
		Quantity:    int(order.Qty.IntPart()),
		Side:        string(order.Side),
		Type:        string(order.Type),
		Status:      string(order.Status),
		FilledQty:   int(order.FilledQty.IntPart()),
		SubmittedAt: order.SubmittedAt,
	}, nil
}

// GetPositions retrieves all positions
func (c *Client) GetPositions() ([]models.Position, error) {
	positions, err := c.alpaca.GetPositions()
	if err != nil {
		return nil, fmt.Errorf("failed to get positions: %w", err)
	}

	result := make([]models.Position, len(positions))
	for i, pos := range positions {
		result[i] = models.Position{
			Symbol:          pos.Symbol,
			Quantity:        int(pos.Qty.IntPart()),
			AvgEntryPrice:   pos.AvgEntryPrice.InexactFloat64(),
			CurrentPrice:    pos.CurrentPrice.InexactFloat64(),
			MarketValue:     pos.MarketValue.InexactFloat64(),
			CostBasis:       pos.CostBasis.InexactFloat64(),
			UnrealizedPL:    pos.UnrealizedPL.InexactFloat64(),
			UnrealizedPLPct: pos.UnrealizedPLPC.InexactFloat64() * 100,
			Side:            string(pos.Side),
		}
	}

	return result, nil
}

// GetOrders retrieves order history
func (c *Client) GetOrders() ([]models.Order, error) {
	orders, err := c.alpaca.GetOrders(alpaca.GetOrdersRequest{
		Status: "all",
		Limit:  100,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}

	result := make([]models.Order, len(orders))
	for i, order := range orders {
		filledAvgPx := 0.0
		if order.FilledAvgPrice != nil {
			filledAvgPx = order.FilledAvgPrice.InexactFloat64()
		}

		result[i] = models.Order{
			ID:          order.ID,
			Symbol:      order.Symbol,
			Quantity:    int(order.Qty.IntPart()),
			Side:        string(order.Side),
			Type:        string(order.Type),
			Status:      string(order.Status),
			FilledQty:   int(order.FilledQty.IntPart()),
			FilledAvgPx: filledAvgPx,
			SubmittedAt: order.SubmittedAt,
			FilledAt:    order.FilledAt,
		}
	}

	return result, nil
}