package conf

import (
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/log"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	EthClient     *ethclient.Client
	RPCURL        string
	StartingBlock uint64
	RPCTimeout    time.Duration
)

func init() {
	rpc, found := os.LookupEnv("RPC_URL")
	if found {
		RPCURL = rpc
	} else {
		RPCURL = "https://eth.llamarpc.com"
	}
	rpc_timeout, found := os.LookupEnv("RPC_TIMEOUT")
	if found {
		parsed, parse_err := strconv.ParseInt(rpc_timeout, 10, 0)
		log.Fatal(parse_err)
		RPCTimeout = time.Duration(parsed)
	} else {
		RPCURL = "https://eth.llamarpc.com"
	}
	st, found := os.LookupEnv("STARTING_BLOCK")

	if found {
		x, err := strconv.ParseInt(st, 10, 64)
		if err != nil {
			log.Fatalf("Network: %s", err)
		}
		StartingBlock = uint64(x)
	} else {
		StartingBlock = 10000
	}
	client, err := ethclient.Dial(RPCURL)
	if err != nil {
		log.Errorf("Client Connection Error : %s  ", err)
	}
	// ctx, _ := context.WithTimeout(context.Background(),time.Second * 10)
	EthClient = client
}

func NetworkValueAddress(chain int64) common.Address {
	return common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE")
}
