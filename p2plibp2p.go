// p2plibp2p.go
package main

import (
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

// HandlePeerFound is called when a new peer is discovered via mDNS.
func (n *Notifee) HandlePeerFound(pi peer.AddrInfo) {
	fmt.Println("Discovered new peer:", pi.ID)

	// Attempt to connect to the discovered peer.
	if err := n.h.Connect(context.Background(), pi); err != nil {
		fmt.Println("Error connecting to peer:", err)
	}
}

// SetupDiscovery starts mDNS discovery service.
func SetupDiscovery(h host.Host) error {
	// Create a new Notifee
	n := &Notifee{h: h}

	// Initialize MDNS service with the notifee
	// Note: NewMdnsService only returns the service, not an error
	_ = mdns.NewMdnsService(h, DiscoveryServiceTag, n)

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
		s.Close()
	})
}

// BroadcastBlockchain sends the current blockchain to all connected peers.
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
		}
		s.Close()
	}
}

// CreateLibp2pHost creates a new libp2p host.
func CreateLibp2pHost() (host.Host, error) {
	// Create a new libp2p host with default options
	h, err := libp2p.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create libp2p host: %w", err)
	}
	return h, nil
}
