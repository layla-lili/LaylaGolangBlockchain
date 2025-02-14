package main



const (
	BlockProtocol = "/block/1.0.0"
	TxProtocol    = "/tx/1.0.0"
	ChainProtocol = "/chain/1.0.0"
)

type Message struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}

type BlockchainDB interface {
	SaveBlock(block Block) error
	GetBlock(hash string) (Block, error)
	SaveChain(chain []Block) error
	LoadChain() ([]Block, error)
}

func SelectBestChain(chains [][]Block) []Block {
	var bestChain []Block
	maxLength := 0

	for _, chain := range chains {
		if len(chain) > maxLength && ValidateBlockchain(chain) {
			maxLength = len(chain)
			bestChain = chain
		}
	}
	return bestChain
}
