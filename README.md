# Terminal Trader

A command-line stock trading application built in Go with real-time market data integration.

## Overview

Terminal Trader provides a CLI interface for paper trading through the Alpaca API. Execute trades, monitor positions, and track portfolio performance directly from your terminal.

## Features

- Real-time stock quotes
- Market order execution (buy/sell)
- Portfolio tracking with P&L calculations
- Order history
- Account management
- Paper trading integration

## Prerequisites

- Go 1.21 or higher
- Alpaca paper trading account (free at alpaca.markets)

## Installation

```bash
git clone https://github.com/vp2306/terminal-trader.git
cd terminal-trader
go mod download
```

## Configuration

Create `config/config.yaml`:

```yaml
alpaca:
  api_key: "YOUR_API_KEY"
  api_secret: "YOUR_SECRET_KEY"
  base_url: "https://paper-api.alpaca.markets"
  data_url: "https://data.alpaca.markets"

trading:
  max_positions: 10
  default_order_size: 10
  risk_per_trade: 0.02
  max_order_value: 10000.0
```

## Usage

```bash
# View account information
go run cmd/terminal-trader/main.go account

# Get stock quote
go run cmd/terminal-trader/main.go quote AAPL

# Buy shares
go run cmd/terminal-trader/main.go buy AAPL -q 10

# Sell shares
go run cmd/terminal-trader/main.go sell AAPL -q 5

# View portfolio
go run cmd/terminal-trader/main.go portfolio

# View order history
go run cmd/terminal-trader/main.go history
```

## Project Structure

```
terminal-trader/
├── cmd/terminal-trader/    # Application entry point
├── internal/
│   ├── api/               # Alpaca API client
│   ├── config/            # Configuration management
│   └── models/            # Data structures
└── config/                # Configuration files
```

## Technology Stack

- Go 1.21+
- Alpaca Trade API v3
- Cobra (CLI framework)
- Viper (configuration)

## Building

```bash
go build -o terminal-trader cmd/terminal-trader/main.go
./terminal-trader account
```
