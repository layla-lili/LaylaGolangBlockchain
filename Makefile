# Variables
PROJECT_ROOT := $(shell pwd)
GO_FILES := $(shell find . -name "*.go" ! -name "all.go" -type f)
TIMESTAMP := $(shell date "+%Y-%m-%d %H:%M:%S")

.PHONY: all consolidate format test clean deps build run

all: consolidate format test

consolidate:
	@echo "🔄 Consolidating files..."
	@{ \
		echo "// This file is auto-generated. Do not edit directly."; \
		echo "// Last updated: $(TIMESTAMP)"; \
		echo ""; \
		echo "package main"; \
		echo ""; \
		echo "import ("; \
		echo '    "bytes"'; \
		echo '    "context"'; \
		echo '    "crypto/ecdsa"'; \
		echo '    "crypto/elliptic"'; \
		echo '    "crypto/rand"'; \
		echo '    "crypto/sha256"'; \
		echo '    "encoding/hex"'; \
		echo '    "encoding/json"'; \
		echo '    "errors"'; \
		echo '    "fmt"'; \
		echo '    "net/http"'; \
		echo '    "os"'; \
		echo '    "strings"'; \
		echo '    "sync"'; \
		echo '    "testing"'; \
		echo '    "time"'; \
		echo ""; \
		echo '    "github.com/cbergoon/merkletree"'; \
		echo '    "github.com/gorilla/mux"'; \
		echo '    "github.com/gorilla/websocket"'; \
		echo '    libp2p "github.com/libp2p/go-libp2p"'; \
		echo '    "github.com/libp2p/go-libp2p/core/host"'; \
		echo '    "github.com/libp2p/go-libp2p/core/network"'; \
		echo '    "github.com/libp2p/go-libp2p/core/peer"'; \
		echo '    mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"'; \
		echo '    "github.com/multiformats/go-multiaddr"'; \
		echo ")"; \
		echo ""; \
	} > all.go
	@if [ -f README.md ]; then \
		echo "/*" >> all.go; \
		cat README.md >> all.go; \
		echo "*/" >> all.go; \
		echo "" >> all.go; \
	fi
	@for file in main.go block.go block_test.go transaction.go chain.go mempool.go merkle.go wallet.go p2p.go p2plibp2p.go server.go consensus.go; do \
		if [ -f $$file ]; then \
			echo "// ======================" >> all.go; \
			echo "// $$file" >> all.go; \
			echo "// ======================" >> all.go; \
			echo "" >> all.go; \
			cat $$file >> all.go; \
			echo "" >> all.go; \
		fi \
	done
	@echo "✅ Files consolidated successfully"

format:
	@echo "🎯 Formatting consolidated file..."
	@go fmt all.go

test:
	@echo "🧪 Running tests..."
	@go test ./...

clean:
	@echo "🧹 Cleaning up..."
	@rm -f all.go