# Makefile for Terminal Trader

.PHONY: all build clean

all: build

build:
	go build -o bin/terminal-trader cmd/terminal-trader/main.go

clean:
	go clean
	rm -f bin/terminal-trader