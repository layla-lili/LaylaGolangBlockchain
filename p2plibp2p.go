// p2plibp2p.go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	// "time"

	libp2p "github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
)

const protocolID = "/blockchain/1.0.0"
const DiscoveryServiceTag = "blockchain-discovery"

// Notifee implements the mdns.Notifee interface for peer discovery.
type Notifee struct {
	h host.Host
}

func (n *Notifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Printf("🎯 Peer Discovered: %s - Attempting connection...\n", pi.ID)

	err := n.h.Connect(context.Background(), pi)
	if err != nil {
		fmt.Printf("❌ Failed to connect to peer %s: %v\n", pi.ID, err)
	} else {
		fmt.Printf("🔗 Successfully connected to peer: %s\n", pi.ID)
	}
}

// SetupDiscovery starts mDNS discovery service.
// func SetupDiscovery(h host.Host) error {
// 	// Create a new Notifee
// 	n := &Notifee{h: h}

// 	// Initialize MDNS service with the notifee
// 	// Note: NewMdnsService only returns the service, not an error
// 	_ = mdns.NewMdnsService(h, DiscoveryServiceTag, n)

//		return nil
//	}
func SetupDiscovery(h host.Host) error {
	n := &Notifee{h: h}

	service := mdns.NewMdnsService(h, DiscoveryServiceTag, n)
	if service == nil {
		return fmt.Errorf("failed to start mDNS service")
	}

	fmt.Println("✅ mDNS discovery service started. Waiting for peers...")

	return nil
}

// SetupStreamHandler registers a handler for incoming streams on our protocol.
func SetupStreamHandler(h host.Host) {
	h.SetStreamHandler(protocolID, func(s network.Stream) {
		fmt.Println("Received new stream from:", s.Conn().RemotePeer().String())
		var receivedChain []Block
		decoder := json.NewDecoder(s)
		if err := decoder.Decode(&receivedChain); err != nil {
			fmt.Println("Error decoding blockchain from stream:", err)
		} else {
			fmt.Println("Received blockchain from peer:")
			for _, b := range receivedChain {
				fmt.Printf("Block %d: %s\n", b.Index, b.Hash)
			}
			// Here, you could add logic to compare and merge chains.
		}

		// Verify received block's Merkle root
		for _, block := range receivedChain {
			merkleRoot, err := GetMerkleRoot(block.Transactions)
			if err != nil {
				fmt.Printf("❌ Failed to create Merkle tree for block %d: %v\n", block.Index, err)
				continue
			}

			if !bytes.Equal(merkleRoot, block.MerkleRoot) {
				fmt.Printf("❌ Invalid Merkle root in block %d\n", block.Index)
				continue
			}

			// Verify all transactions in the block
			for _, tx := range block.Transactions {
				isValid, err := VerifyTransactionInBlock(block, tx)
				if err != nil {
					fmt.Printf("⚠️ Warning: Failed to verify transaction %s: %v\n", tx.TxID, err)
					continue
				}

				if !isValid {
					fmt.Printf("❌ Invalid transaction detected in block %d: %s\n", block.Index, tx.TxID)
					// Handle invalid transaction (maybe reject the block)
				}
			}
		}

		s.Close()
	})
}

// BroadcastBlockchain sends the current blockchain to all connected peers.
//
//	func BroadcastBlockchain(h host.Host, blockchain []Block) {
//		peers := h.Network().Peers()
//		fmt.Println("Broadcasting blockchain to", len(peers), "peers.")
//		for _, p := range peers {
//			s, err := h.NewStream(context.Background(), p, protocolID)
//			if err != nil {
//				fmt.Println("Error opening stream to peer", p.String(), ":", err)
//				continue
//			}
//			encoder := json.NewEncoder(s)
//			if err := encoder.Encode(blockchain); err != nil {
//				fmt.Println("Error sending blockchain to peer", p.String(), ":", err)
//			}
//			s.Close()
//		}
//	}
func BroadcastBlockchain(h host.Host, blockchain []Block) {
	peers := h.Network().Peers()
	fmt.Println("Broadcasting blockchain to", len(peers), "peers.")
	for _, p := range peers {
		s, err := h.NewStream(context.Background(), p, protocolID)
		if err != nil {
			fmt.Println("Error opening stream to peer", p.String(), ":", err)
			continue
		}
		encoder := json.NewEncoder(s)
		if err := encoder.Encode(blockchain); err != nil {
			fmt.Println("Error sending blockchain to peer", p.String(), ":", err)
		} else {
			fmt.Printf("Blockchain sent to peer %s\n", p.String())
		}
		s.Close()
	}
}

// CreateLibp2pHost creates a new libp2p host.
func CreateLibp2pHost() (host.Host, error) {
	// Create a new libp2p host with default options
	fmt.Println("🚀 Creating libp2p host...")
	h, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}
	return h, nil
}
