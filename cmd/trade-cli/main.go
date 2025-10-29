package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vp2306/TradeCLI/internal/api"
	"github.com/vp2306/TradeCLI/internal/config"
)

var (
	cfg    *config.Config
	client *api.Client
)

func main() {
	var err error
	cfg, err = config.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		fmt.Println("Make sure config/config.yaml exists with your Alpaca API keys")
		os.Exit(1)
	}

	client = api.NewClient(cfg)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "terminal-trader",
	Short: "Terminal Trader - CLI Stock Trading Application",
	Long:  `A terminal-based stock trading engine with real-time market data and order execution.`,
}

var accountCmd = &cobra.Command{
	Use:   "account",
	Short: "Display account information",
	Run: func(cmd *cobra.Command, args []string) {
		acct, err := client.GetAccount()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Println("\n=== Account Information ===")
		fmt.Printf("Account ID:         %s\n", acct.ID)
		fmt.Printf("Cash:               $%.2f\n", acct.Cash)
		fmt.Printf("Portfolio Value:    $%.2f\n", acct.PortfolioValue)
		fmt.Printf("Buying Power:       $%.2f\n", acct.BuyingPower)
		fmt.Printf("Equity:             $%.2f\n", acct.Equity)
		fmt.Printf("Daytrade Count:     %d\n", acct.DaytradeCount)
		fmt.Printf("Pattern Day Trader: %v\n", acct.PatternDayTrader)
		fmt.Println()
	},
}

var quoteCmd = &cobra.Command{
	Use:   "quote [SYMBOL]",
	Short: "Get real-time quote for a stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		symbol := strings.ToUpper(args[0])
		quote, err := client.GetQuote(symbol)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		fmt.Printf("\n=== Quote: %s ===\n", quote.Symbol)
		fmt.Printf("Last Price: $%.2f\n", quote.LastPrice)
		fmt.Printf("Volume:     %d\n", quote.Volume)
		fmt.Printf("Timestamp:  %s\n", quote.Timestamp.Format("2006-01-02 15:04:05"))
		fmt.Println()
	},
}

var buyCmd = &cobra.Command{
	Use:   "buy [SYMBOL]",
	Short: "Buy shares of a stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		symbol := strings.ToUpper(args[0])
		qty, _ := cmd.Flags().GetInt("quantity")

		if qty <= 0 {
			fmt.Println("Error: quantity must be greater than 0")
			return
		}

		order, err := client.PlaceOrder(symbol, qty, "buy")
		if err != nil {
			fmt.Printf("Error placing order: %v\n", err)
			return
		}

		fmt.Printf("\n✓ Buy order placed successfully!\n")
		fmt.Printf("Order ID:  %s\n", order.ID)
		fmt.Printf("Symbol:    %s\n", order.Symbol)
		fmt.Printf("Quantity:  %d\n", order.Quantity)
		fmt.Printf("Status:    %s\n", order.Status)
		fmt.Println()
	},
}

var sellCmd = &cobra.Command{
	Use:   "sell [SYMBOL]",
	Short: "Sell shares of a stock",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		symbol := strings.ToUpper(args[0])
		qty, _ := cmd.Flags().GetInt("quantity")

		if qty <= 0 {
			fmt.Println("Error: quantity must be greater than 0")
			return
		}

		order, err := client.PlaceOrder(symbol, qty, "sell")
		if err != nil {
			fmt.Printf("Error placing order: %v\n", err)
			return
		}

		fmt.Printf("\n✓ Sell order placed successfully!\n")
		fmt.Printf("Order ID:  %s\n", order.ID)
		fmt.Printf("Symbol:    %s\n", order.Symbol)
		fmt.Printf("Quantity:  %d\n", order.Quantity)
		fmt.Printf("Status:    %s\n", order.Status)
		fmt.Println()
	},
}

var portfolioCmd = &cobra.Command{
	Use:   "portfolio",
	Short: "Display current portfolio positions",
	Run: func(cmd *cobra.Command, args []string) {
		positions, err := client.GetPositions()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(positions) == 0 {
			fmt.Println("\nNo open positions.\n")
			return
		}

		fmt.Println("\n=== Portfolio ===")
		fmt.Printf("%-8s %-6s %-12s %-12s %-12s %-12s %-10s\n",
			"Symbol", "Qty", "Avg Entry", "Current", "Market Val", "P&L", "P&L %")
		fmt.Println(strings.Repeat("-", 80))

		totalPL := 0.0
		for _, pos := range positions {
			plSign := "+"
			if pos.UnrealizedPL < 0 {
				plSign = ""
			}

			fmt.Printf("%-8s %-6d $%-11.2f $%-11.2f $%-11.2f %s$%-10.2f %s%.2f%%\n",
				pos.Symbol,
				pos.Quantity,
				pos.AvgEntryPrice,
				pos.CurrentPrice,
				pos.MarketValue,
				plSign,
				pos.UnrealizedPL,
				plSign,
				pos.UnrealizedPLPct)

			totalPL += pos.UnrealizedPL
		}

		fmt.Println(strings.Repeat("-", 80))
		plSign := "+"
		if totalPL < 0 {
			plSign = ""
		}
		fmt.Printf("Total Unrealized P&L: %s$%.2f\n\n", plSign, totalPL)
	},
}

var historyCmd = &cobra.Command{
	Use:   "history",
	Short: "Display order history",
	Run: func(cmd *cobra.Command, args []string) {
		orders, err := client.GetOrders()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}

		if len(orders) == 0 {
			fmt.Println("\nNo order history.\n")
			return
		}

		fmt.Println("\n=== Order History ===")
		fmt.Printf("%-8s %-6s %-4s %-6s %-10s %-12s %-19s\n",
			"Symbol", "Side", "Qty", "Type", "Status", "Filled Px", "Submitted At")
		fmt.Println(strings.Repeat("-", 85))

		for _, order := range orders {
			filledPx := "N/A"
			if order.FilledAvgPx > 0 {
				filledPx = fmt.Sprintf("$%.2f", order.FilledAvgPx)
			}

			fmt.Printf("%-8s %-6s %-4d %-6s %-10s %-12s %s\n",
				order.Symbol,
				order.Side,
				order.Quantity,
				order.Type,
				order.Status,
				filledPx,
				order.SubmittedAt.Format("2006-01-02 15:04:05"))
		}
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(accountCmd)
	rootCmd.AddCommand(quoteCmd)
	rootCmd.AddCommand(buyCmd)
	rootCmd.AddCommand(sellCmd)
	rootCmd.AddCommand(portfolioCmd)
	rootCmd.AddCommand(historyCmd)

	buyCmd.Flags().IntP("quantity", "q", 1, "Number of shares to buy")
	sellCmd.Flags().IntP("quantity", "q", 1, "Number of shares to sell")
}