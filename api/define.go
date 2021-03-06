package api

import (
	"math/big"

	"github.com/seeleteam/go-seele/common"
	"github.com/seeleteam/go-seele/core/state"
	"github.com/seeleteam/go-seele/core/store"
	"github.com/seeleteam/go-seele/core/types"
	"github.com/seeleteam/go-seele/database"
	"github.com/seeleteam/go-seele/log"
	"github.com/seeleteam/go-seele/p2p"
	"github.com/seeleteam/go-seele/rpc"
)

// Backend interface provides the common API services (that are provided by
// both full and light clients) with access to necessary functions.
type Backend interface {
	GetP2pServer() *p2p.Server
	GetNetVersion() string
	GetNetWorkID() string

	TxPoolBackend() Pool
	ChainBackend() Chain
	ProtocolBackend() Protocol
	Log() *log.SeeleLog
	IsSyncing() bool

	GetAccountIndexDB() database.Database
	GetIndexAccountDB() database.Database

	GetBlock(hash common.Hash, height int64) (*types.Block, error)
	GetBlockTotalDifficulty(hash common.Hash) (*big.Int, error)
	GetReceiptByTxHash(txHash common.Hash) (*types.Receipt, error)
	GetTransaction(pool PoolCore, bcStore store.BlockchainStore, txHash common.Hash) (*types.Transaction, *BlockIndex, error)

	// // // Calculate the proposer
	// // CalcProposer(lastProposer common.Address, round uint64)
	// // // Return the Verifier size
	// // Size() int
	// // // Return the Verifier array
	// List() []bft.Verifier
	// // // Get Verifier by index
	// // GetVerByIndex(i uint64) bft.Verifier
	// // // Get Verifier by given address
	// // GetVerByAddress(addr common.Address) (int, bft.Verifier)
	// // // Get current proposer
	// GetProposer() bft.Verifier
	// // // Check whether the Verifier with given address is a proposer
	// // IsProposer(address common.Address) bool
	// // // Add Verifier
	// // AddVerifier(address common.Address) bool
	// // // Remove Verifier
	// // RemoveVerifier(address common.Address) bool
	// // // Copy Verifier set
	// // Copy() bft.VerifierSet
	// // // Get the maximum number of faulty nodes
	// // F() int
	// // // Get proposer policy
	// // Policy() bft.ProposerPolicy
}

func GetAPIs(apiBackend Backend) []rpc.API {
	return []rpc.API{
		{
			Namespace: "seele",
			Version:   "1.0",
			Service:   NewPublicSeeleAPI(apiBackend),
			Public:    true,
		},
		{
			Namespace: "txpool",
			Version:   "1.0",
			Service:   NewTransactionPoolAPI(apiBackend),
			Public:    true,
		},
		{
			Namespace: "network",
			Version:   "1.0",
			Service:   NewPrivateNetworkAPI(apiBackend),
			Public:    false,
		},
		{
			Namespace: "debug",
			Version:   "1.0",
			Service:   NewPrivateDebugAPI(apiBackend),
			Public:    false,
		},
		{
			Namespace: "subchain",
			Version:   "1.0",
			Service:   NewPublicSubchainAPI(apiBackend),
			Public:    true,
		}}
}

// MinerInfo miner simple info
type GetMinerInfo struct {
	Coinbase           common.Address
	CurrentBlockHeight uint64
	HeaderHash         common.Hash
	Shard              uint
	MinerStatus        string
	Version            string
	BlockAge           *big.Int
	PeerCnt            string
}

// GetBalanceResponse response param for GetBalance api
type GetBalanceResponse struct {
	Account common.Address
	Balance *big.Int
}

// GetLogsResponse response param for GetLogs api
type GetLogsResponse struct {
	*types.Log
	Txhash   common.Hash
	LogIndex uint
	Args     interface{} `json:"data"`
}

type PoolCore interface {
	AddTransaction(tx *types.Transaction) error
	GetTransaction(txHash common.Hash) *types.Transaction
}

type Pool interface {
	PoolCore
	GetTransactions(processing, pending bool) []*types.Transaction
	GetTxCount() int
}

type Chain interface {
	CurrentHeader() *types.BlockHeader
	GetCurrentState() (*state.Statedb, error)
	GetState(blockHash common.Hash) (*state.Statedb, error)
	GetStore() store.BlockchainStore
}

type Protocol interface {
	SendDifferentShardTx(tx *types.Transaction, shard uint)
	GetProtocolVersion() (uint, error)
}

type BFTCore interface {
}
