// p2plibp2p.go
package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"time"

	libp2p "github.com/libp2p/go-libp2p"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	mdns "github.com/libp2p/go-libp2p/p2p/discovery/mdns"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/multiformats/go-multiaddr"
)

const protocolID = "/blockchain/1.0.0"
const DiscoveryServiceTag = "blockchain-discovery"

// Notifee implements the mdns.Notifee interface for peer discovery.
type Notifee struct {
	h host.Host
}

func (n *Notifee) HandlePeerFound(pi peer.AddrInfo) {
	if pi.ID == n.h.ID() {
		return // Skip self-connection
	}

	fmt.Printf("🎯 Peer Discovered: %s\n", pi.ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := n.h.Connect(ctx, pi); err != nil {
		fmt.Printf("❌ Failed to connect to peer %s: %v\n", pi.ID, err)
		return
	}

	fmt.Printf("🔗 Successfully connected to peer: %s\n", pi.ID)

	// Open stream with retry
	var stream network.Stream
	var err error
	for i := 0; i < 3; i++ {
		stream, err = n.h.NewStream(ctx, pi.ID, protocolID)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
	}
	if err != nil {
		fmt.Printf("❌ Failed to open stream after retries: %v\n", err)
		return
	}
	defer stream.Close()

	// Send hello message
	message := fmt.Sprintf("Hello from %s!", n.h.ID())
	if _, err := stream.Write([]byte(message)); err != nil {
		fmt.Printf("❌ Failed to send hello: %v\n", err)
		return
	}
}

// Add helper methods for stream handling
func (n *Notifee) tryOpenStream(peerID peer.ID) (network.Stream, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	stream, err := n.h.NewStream(ctx, peerID, protocolID)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %w", err)
	}
	return stream, nil
}

func (n *Notifee) sendHello(stream network.Stream) error {
	message := fmt.Sprintf("Hello from %s!", n.h.ID().String())
	_, err := stream.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to send hello: %w", err)
	}
	return nil
}

// SetupDiscovery starts mDNS discovery service.
func SetupDiscovery(h host.Host) error {
	// Create a new discovery service
	discovery := mdns.NewMdnsService(
		h,
		DiscoveryServiceTag,
		&Notifee{h: h},
	)

	if discovery == nil {
		return fmt.Errorf("failed to create discovery service")
	}

	// Start the discovery service
	if err := discovery.Start(); err != nil {
		return fmt.Errorf("failed to start discovery service: %w", err)
	}

	fmt.Printf("✅ mDNS discovery started with tag: %s\n", DiscoveryServiceTag)
	return nil
}

// SetupStreamHandler registers a handler for incoming streams on our protocol.
func SetupStreamHandler(h host.Host, state *BlockchainState) {
	h.SetStreamHandler(protocolID, func(s network.Stream) {
		// Read the hello message first
		buf := make([]byte, 1024)
		n, err := s.Read(buf)
		if err != nil {
			fmt.Printf("❌ Error reading hello message: %v\n", err)
			s.Reset()
			return
		}

		helloMsg := string(buf[:n])
		fmt.Printf("📨 Received: %s\n", helloMsg)

		// Now handle blockchain sync
		var receivedChain []Block
		decoder := json.NewDecoder(s)
		if err := decoder.Decode(&receivedChain); err != nil {
			// This might not be a chain sync message, that's okay
			return
		}

		consensus := state.GetConsensus()
		if consensus.HandleChainSync(receivedChain) {
			fmt.Println("✅ Chain updated successfully")
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
func CreateLibp2pHost(port string) (host.Host, error) {
	// Generate new private key
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.Ed25519, 2048, rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key pair: %w", err)
	}

	// Create multiaddress
	addr := fmt.Sprintf("/ip4/0.0.0.0/tcp/%s", port)
	ma, err := multiaddr.NewMultiaddr(addr)
	if err != nil {
		return nil, fmt.Errorf("failed to create multiaddr: %w", err)
	}

	// Create connection manager first
	connManager, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection manager: %w", err)
	}

	// Create libp2p host with the connection manager
	h, err := libp2p.New(
		libp2p.ListenAddrs(ma),
		libp2p.Identity(priv),
		libp2p.DefaultTransports,
		libp2p.DefaultMuxers,
		libp2p.DefaultSecurity,
		libp2p.NATPortMap(),
		libp2p.DisableRelay(),
		libp2p.ConnectionManager(connManager),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create host: %w", err)
	}

	// Log the host's addresses
	fmt.Printf("🌍 P2P Host created with ID: %s\n", h.ID().String())
	fmt.Println("📡 Listening on addresses:")
	for _, addr := range h.Addrs() {
		fmt.Printf("   - %s/p2p/%s\n", addr, h.ID().String())
	}

	return h, nil
}
