# Blockchain Implementation in Go

## Table of Contents
- [Blockchain Implementation in Go](#blockchain-implementation-in-go)
  - [Table of Contents](#table-of-contents)
- [Blockchain Implementation in Go](#blockchain-implementation-in-go-1)
  - [Code Flow Explanation](#code-flow-explanation)
    - [1. Application Initialization (`main.go`)](#1-application-initialization-maingo)
      - [State Initialization:](#state-initialization)
      - [Genesis Block Creation:](#genesis-block-creation)
      - [Wallet Initialization:](#wallet-initialization)
      - [P2P Host Setup:](#p2p-host-setup)
      - [Server Initialization:](#server-initialization)
    - [2. Block Handling (`block.go` \& `block_test.go`)](#2-block-handling-blockgo--block_testgo)
      - [Block Structure:](#block-structure)
      - [Hash Calculation:](#hash-calculation)
      - [Block Generation \& Mining:](#block-generation--mining)
      - [Genesis Block:](#genesis-block)
      - [Testing:](#testing)
    - [3. Transaction Processing (`transaction.go`)](#3-transaction-processing-transactiongo)
      - [Transaction Structure:](#transaction-structure)
      - [TxID Calculation:](#txid-calculation)
      - [Signature and Verification:](#signature-and-verification)
    - [4. Merkle Tree Operations (`merkle.go`)](#4-merkle-tree-operations-merklego)
      - [Integration with External Library:](#integration-with-external-library)
      - [Merkle Tree Interface:](#merkle-tree-interface)
      - [Merkle Root Calculation \& Verification:](#merkle-root-calculation--verification)
    - [5. Blockchain State and Server (`server.go`)](#5-blockchain-state-and-server-servergo)
      - [Server Object:](#server-object)
      - [API Endpoints:](#api-endpoints)
      - [Middleware:](#middleware)
    - [6. P2P Communication (`p2plibp2p.go` \& `p2p.go`)](#6-p2p-communication-p2plibp2pgo--p2pgo)
      - [Discovery \& Connection:](#discovery--connection)
      - [Stream Handling:](#stream-handling)
      - [Broadcasting:](#broadcasting)
    - [7. Consensus (`consensus.go`)](#7-consensus-consensusgo)
      - [Chain Synchronization:](#chain-synchronization)
  - [CLI Interface](#cli-interface)
    - [CLI Commands](#cli-commands)
    - [Running Multiple Nodes](#running-multiple-nodes)
    - [Testing Network Synchronization](#testing-network-synchronization)
  - [Updated Components](#updated-components)
    - [Connection Manager](#connection-manager)
    - [Improved Transaction Validation](#improved-transaction-validation)
    - [Enhanced P2P Communication](#enhanced-p2p-communication)
    - [State Management](#state-management)
  - [Network Scenarios](#network-scenarios)
    - [Adding a New Node](#adding-a-new-node)
    - [Creating Transactions Between Nodes](#creating-transactions-between-nodes)
    - [Handling Network Partitions](#handling-network-partitions)
  - [Common Issues and Solutions](#common-issues-and-solutions)
    - [Peer Discovery Issues](#peer-discovery-issues)
    - [Transaction Failures](#transaction-failures)
    - [Mining Issues](#mining-issues)
  - [Performance Considerations](#performance-considerations)
  - [Evaluation](#evaluation)
    - [Strengths:](#strengths)
      - [Modular Structure:](#modular-structure)
      - [State Encapsulation:](#state-encapsulation)
      - [Integrated P2P Networking:](#integrated-p2p-networking)
      - [Testing Coverage:](#testing-coverage)
      - [Clear API Endpoints:](#clear-api-endpoints)
    - [Future Improvements:](#future-improvements)
  - [Features](#features)
    - [Currently Implemented](#currently-implemented)
    - [Block Features](#block-features)
    - [Transaction Features](#transaction-features)
    - [Network Features](#network-features)
    - [Storage and State](#storage-and-state)
    - [Security](#security)
    - [Planned Future Features](#planned-future-features)
  - [License](#license)

# Blockchain Implementation in Go

This repository contains a basic implementation of a blockchain using Go. The system incorporates essential blockchain components like blocks, transactions, Merkle Trees, and a P2P network. Below is a detailed explanation of the code flow, components, and functionality.

## Code Flow Explanation

### 1. Application Initialization (`main.go`)

#### State Initialization:
The program begins in `main()`, where a new instance of a centralized state (`BlockchainState`) is created via `NewBlockchainState()`. This state encapsulates all key components such as the blockchain (chain of blocks), the wallet, the pending transactions (mempool), and the P2P host.

#### Genesis Block Creation:
The genesis block is created by calling `CreateGenesisBlock()`. This block serves as the starting point of the blockchain. It is then added to the state with `state.AddBlock(genesisBlock)`.

#### Wallet Initialization:
A new wallet is created by calling `NewWallet()`. The wallet (which contains the private/public keys and a derived address) is stored in the state using `state.SetWallet(wallet)`.

#### P2P Host Setup:
A libp2p host is created with `CreateLibp2pHost()`. This host enables the node to participate in a peer-to-peer network. The host is saved to the state via `state.SetP2PHost(p2pHost)`.

Then, the application sets up P2P discovery (using mDNS) with `SetupDiscovery(p2pHost)` and configures the stream handler via `SetupStreamHandler(p2pHost)`. These functions allow the node to find and connect to peers and exchange blockchain data.

#### Server Initialization:
A new HTTP server is created by calling `NewServer(state)`. This server uses the centralized state to serve API endpoints. The server then starts listening on the specified port (default "8080") using `server.Start(apiPort)`.

### 2. Block Handling (`block.go` & `block_test.go`)

#### Block Structure:
The `Block` struct holds an index, timestamp, list of transactions, previous block hash, current block hash, nonce (for mining), Merkle root of its transactions, and a difficulty level.

#### Hash Calculation:
The function `CalculateBlockHash(block Block)` constructs a record string using key fields (including the Merkle root and difficulty) and computes its SHA‑256 hash. This hash uniquely identifies the block.

#### Block Generation & Mining:
- `GenerateBlock(prevBlock, transactions)` creates a new block by incrementing the previous block’s index and setting up the new block’s fields.
- The Merkle root is computed from the transactions via `GetMerkleRoot(transactions)`.
- The block is then “mined” by iterating (incrementing the nonce) until the computed hash has the required number of leading zeros (as defined by the block’s difficulty).

#### Genesis Block:
The genesis block is created in a simplified manner by hashing a string that includes the word "Genesis" along with the timestamp and a nonce.

#### Testing:
The file `block_test.go` contains unit tests that verify:
- The correctness of the Merkle root computation.
- Block generation (index, previous hash linkage, and overall blockchain validity).
- That the mining process produces a hash with the correct difficulty prefix.

### 3. Transaction Processing (`transaction.go`)

#### Transaction Structure:
A `Transaction` holds information such as the sender, receiver, amount, a computed transaction ID (TxID), a signature, timestamp, and fee.

#### TxID Calculation:
`CalculateTxID(tx Transaction)` computes a SHA‑256 hash over a concatenation of the sender, receiver, and amount. This is used as the unique identifier for the transaction.

#### Signature and Verification:
- The wallet’s `SignTransaction` method signs the transaction’s TxID using ECDSA (with ASN.1 encoding).
- `ValidateTransaction` uses the public key (provided as a byte slice) to verify the transaction’s signature.
  
**Note:** There is an expectation that the public key provided for validation is the full key, not merely the derived address.

### 4. Merkle Tree Operations (`merkle.go`)

#### Integration with External Library:
The code integrates with the `github.com/cbergoon/merkletree` package.

#### Merkle Tree Interface:
The `Transaction` type implements the required methods `CalculateHash()` and `Equals()` so that it can be used as content in the Merkle tree.

#### Merkle Root Calculation & Verification:
Functions such as `GetMerkleRoot()` build the tree from a slice of transactions, and helper functions like `VerifyTransactionInBlock()` check if a given transaction is part of the block's Merkle tree.

### 5. Blockchain State and Server (`server.go`)

#### Server Object:
A `Server` struct is defined that holds a pointer to the `BlockchainState`.

#### API Endpoints:
The server registers several endpoints:
- `GET /chain`: Returns the current blockchain.
- `POST /transaction`: Accepts a new transaction. It decodes the transaction JSON, generates a TxID if needed, signs the transaction using the wallet stored in the state, and then adds it to the pending transactions.
- `GET /mine`: Retrieves pending transactions, creates a new block using `GenerateBlock()`, adds it to the chain, and broadcasts the updated chain to peers.
- `GET /peers`: Returns a list of currently connected P2P peers.

#### Middleware:
A simple logging middleware prints out each incoming request.

### 6. P2P Communication (`p2plibp2p.go` & `p2p.go`)

#### Discovery & Connection:
The code uses libp2p along with mDNS for discovering peers. The `Notifee` struct is implemented to handle newly discovered peers by attempting to connect to them.

#### Stream Handling:
The `SetupStreamHandler` function registers a handler for the custom protocol (`"/blockchain/1.0.0"`). When a new stream is received, it decodes a blockchain from the peer, verifies each block’s Merkle root, and validates each transaction.

#### Broadcasting:
The `BroadcastBlockchain` function sends the current blockchain to all connected peers by opening new streams and encoding the chain as JSON.

### 7. Consensus (`consensus.go`)

#### Chain Synchronization:
A simple consensus mechanism is implemented via the `Consensus` struct.

`HandleChainSync(receivedChain []Block)` compares a received chain with the local one: it validates the chain’s integrity, ensures the new chain is longer, and that every block meets the proof-of-work requirements.

Helper functions like `ValidateChain`, `ValidateBlockTransactions`, and `ValidateProofOfWork` are used to enforce these rules.

## CLI Interface

The blockchain node comes with an interactive CLI interface for easy interaction with the network.

### CLI Commands

```bash
# Start the CLI interface automatically when running a node
go run . -http 8080 -p2p 6001

💻 Starting CLI interface...

🚀 Blockchain CLI
1. View blockchain - Shows all blocks in the chain
2. Create transaction - Create a new transaction
3. Mine block - Mine pending transactions into a new block
4. View peers - List all connected P2P nodes
5. Exit
```

### Running Multiple Nodes

Example of running a 4-node network:

```bash
# Terminal 1 - First Node (Bootstrap)
go run . -http 8080 -p2p 6001

# Terminal 2 - Second Node
go run . -http 8081 -p2p 6002

# Terminal 3 - Third Node
go run . -http 8082 -p2p 6003

# Terminal 4 - Fourth Node
go run . -http 8083 -p2p 6004
```

### Testing Network Synchronization

1. Create a transaction on Node 1:
   ```bash
   # In Node 1's CLI
   Choose option 2 (Create transaction)
   ```

2. Mine the block on Node 1:
   ```bash
   # In Node 1's CLI
   Choose option 3 (Mine block)
   ```

3. Verify synchronization on other nodes:
   ```bash
   # In other nodes' CLI
   Choose option 1 (View blockchain)
   ```

## Updated Components

### Connection Manager

The P2P network now includes a connection manager with the following features:
```go
connManager, err := connmgr.NewConnManager(
    100,    // LowWater - minimum number of connections to maintain
    400,    // HighWater - maximum number of connections before pruning
    connmgr.WithGracePeriod(time.Minute), // Grace period for new connections
)
```

### Improved Transaction Validation

Transaction validation now includes multiple checks:
- Public key verification
- Signature validation
- Address derivation verification
- Timestamp validation

### Enhanced P2P Communication

The P2P protocol now includes:
- Automatic peer discovery using mDNS
- Connection retry mechanism
- Stream-based communication
- Blockchain synchronization
- Transaction broadcasting

### State Management

The `BlockchainState` struct manages:
- Blockchain data
- Pending transactions
- Mempool
- Wallet
- P2P connections
- Consensus rules

## Network Scenarios

### Adding a New Node
```bash
# Start a new node with unique ports
go run . -http 8084 -p2p 6005

# The node will:
1. Generate a new wallet
2. Create genesis block
3. Discover peers via mDNS
4. Sync blockchain automatically
```

### Creating Transactions Between Nodes
1. Get recipient's address from Node 2:
   ```bash
   # In Node 2's CLI
   Choose option 4 (View peers)
   ```

2. Create transaction on Node 1:
   ```bash
   # In Node 1's CLI
   Choose option 2 (Create transaction)
   Enter recipient's address
   Enter amount
   ```

3. Mine the transaction:
   ```bash
   # Can be done on any node
   Choose option 3 (Mine block)
   ```

### Handling Network Partitions

The system automatically:
- Retries failed connections
- Maintains peer list
- Resynchronizes when connection is restored
- Validates and resolves chain conflicts

## Common Issues and Solutions

### Peer Discovery Issues
- Ensure mDNS is not blocked by firewall
- Check all nodes are on the same network
- Verify port availability

### Transaction Failures
- Verify sender has sufficient balance
- Check transaction signature
- Ensure receiver address is valid

### Mining Issues
- Check node synchronization
- Verify pending transactions
- Ensure proper difficulty level

## Performance Considerations

- **Memory Usage**: Mempool cleanup occurs every hour
- **Network Bandwidth**: Blockchain syncs only transfer missing blocks
- **CPU Usage**: Mining difficulty adjusts based on network hash rate

## Evaluation

### Strengths:

#### Modular Structure:
The code is divided into multiple files according to functionality (blocks, transactions, state, P2P, server, consensus). This modular design improves readability and maintainability.

#### State Encapsulation:
With the introduction of a `BlockchainState` (even though its implementation isn’t shown here), there’s an effort to encapsulate global state. This is beneficial for testability and future scalability.

#### Integrated P2P Networking:
The use of libp2p and mDNS for peer discovery and blockchain broadcasting is a strong point. The code not only sets up the P2P host but also handles incoming streams and broadcasts new chain data.

#### Testing Coverage:
Unit tests in `block_test.go` provide a good starting point for verifying block generation, mining, and Merkle tree operations.

#### Clear API Endpoints:
The HTTP API exposes key operations (getting the chain, submitting transactions, mining, listing peers) which makes it easier to interact with the blockchain for further development or integration with other systems.

---

### Future Improvements:
- **Smart Contracts**: Introduce a simple smart contract system for token transfer or other decentralized applications.
- **Proof of Stake (PoS)**: Implement a consensus mechanism based on staking.
- **Private Transactions**: Implement transaction privacy features like ZK-SNARKs or ring signatures.

---

## Features

### Currently Implemented
- ✅ Proof of Work consensus mechanism
- ✅ P2P networking using libp2p with mDNS discovery
- ✅ Transaction mempool for managing pending transactions
- ✅ Native wallet implementation with ECDSA key pairs
- ✅ Merkle tree for transaction verification
- ✅ RESTful API for blockchain interaction
- ✅ Interactive CLI interface
- ✅ Block validation with difficulty adjustment
- ✅ Transaction signing and verification
- ✅ Multi-node synchronization

### Block Features
- ✅ Dynamic difficulty adjustment
- ✅ Merkle root for transaction verification
- ✅ Block headers with timestamps
- ✅ Previous block hash linking
- ❌ Block size limits
- ❌ Chain reorganization for forks

### Transaction Features
- ✅ ECDSA transaction signing
- ✅ Transaction pool management
- ✅ Basic transaction validation
- ❌ Transaction batching
- ❌ UTXO model
- ❌ Multi-signature support

### Network Features
- ✅ Peer discovery via mDNS
- ✅ Blockchain synchronization
- ✅ Transaction broadcasting
- ✅ Connection management
- ❌ Chain history pruning
- ❌ IPFS integration

### Storage and State
- ✅ In-memory blockchain state
- ✅ Transaction mempool
- ✅ Peer connection state
- ❌ Persistent storage
- ❌ State checkpoints
- ❌ Immutable action logs

### Security
- ✅ ECDSA key pair generation
- ✅ Transaction signature verification
- ✅ Basic address generation
- ❌ Token staking
- ❌ Governance mechanisms
- ❌ Advanced cryptographic schemes

### Planned Future Features
1. Chain reorganization to handle forks
2. Block size and transaction limits
3. Transaction batching for improved throughput
4. Chain history pruning
5. Persistent storage with database integration
6. IPFS integration for distributed file storage
7. Token staking and governance mechanisms
8. Immutable audit logs

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
