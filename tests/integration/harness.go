//go:build integration

package integration

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	geth_common "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/1inch/1inch-sdk-go/v4/common"
	"github.com/1inch/1inch-sdk-go/v4/constants"
	web3_provider "github.com/1inch/1inch-sdk-go/v4/internal/web3-provider"
)

const (
	wethAddress      = "0xc02aaa39b223fe8d0a0e5c4f27ead9083c756cc2"
	usdcAddress      = "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"
	accessToken      = "0xacce550000159e70908c0499a1119d04e7039c28"
	usdcDonorAddress = "0x37305B1cD40574E4C5Ce33f8e8306Be057fD7341"

	lopV4Address = constants.AggregationRouterV6
	zeroAddress  = constants.ZeroAddress
)

// freshPrivateKey generates a random key so test accounts have no mainnet state.
// Well-known keys (e.g. anvil defaults) must be avoided: their mainnet accounts carry
// EIP-7702 delegations, which gives them code and breaks Permit2's ecrecover path.
func freshPrivateKey(t *testing.T) string {
	t.Helper()
	key, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	return fmt.Sprintf("%x", crypto.FromECDSA(key))
}

var fallbackForkUrls = []string{
	"https://ethereum-rpc.publicnode.com",
	"https://eth.merkle.io",
	"https://1rpc.io/eth",
	"https://eth.drpc.org", // free tier throttles quickly, last resort
}

// anvilCmd holds the shared fork process, torn down in TestMain after all tests run
var anvilCmd *exec.Cmd

// TestMain tears down the shared anvil fork once the whole package has run
func TestMain(m *testing.M) {
	code := m.Run()
	if anvilCmd != nil && anvilCmd.Process != nil {
		_ = anvilCmd.Process.Kill()
		_, _ = anvilCmd.Process.Wait()
	}
	os.Exit(code)
}

type forkNode struct {
	url       string
	rpcClient *rpc.Client
	ethClient *ethclient.Client
}

func resolveForkUrl(t *testing.T) string {
	t.Helper()
	if url := os.Getenv("FORK_URL"); url != "" {
		return url
	}
	for _, url := range fallbackForkUrls {
		client, err := rpc.Dial(url)
		if err != nil {
			continue
		}
		var blockNumber string
		err = client.CallContext(context.Background(), &blockNumber, "eth_blockNumber")
		client.Close()
		if err == nil && blockNumber != "" {
			t.Logf("using fallback fork url: %s", url)
			return url
		}
	}
	t.Fatal("no working fork url: set FORK_URL to a mainnet RPC endpoint")
	return ""
}

func freePort(t *testing.T) int {
	t.Helper()
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("failed to find free port: %v", err)
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port
}

// startAnvil boots an anvil mainnet fork as a subprocess and returns RPC handles to it
func startAnvil(t *testing.T) *forkNode {
	t.Helper()

	if _, err := exec.LookPath("anvil"); err != nil {
		t.Fatal("anvil not found on PATH: install foundry (https://getfoundry.sh)")
	}

	forkUrl := resolveForkUrl(t)
	port := freePort(t)

	cmd := exec.Command("anvil",
		"-f", forkUrl,
		"--chain-id", "1",
		"--port", fmt.Sprintf("%d", port),
	)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		t.Fatalf("failed to get anvil stdout: %v", err)
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		t.Fatalf("failed to start anvil: %v", err)
	}
	anvilCmd = cmd

	ready := make(chan struct{})
	died := make(chan string, 1)
	go func() {
		scanner := bufio.NewScanner(stdout)
		signaled := false
		lastLine := ""
		// Keep draining to EOF for the process lifetime so anvil never blocks on a
		// full pipe; EOF before the readiness line means the process exited early
		for scanner.Scan() {
			lastLine = scanner.Text()
			if !signaled && strings.Contains(lastLine, "Listening on") {
				signaled = true
				close(ready)
			}
		}
		if !signaled {
			died <- lastLine
		}
	}()

	select {
	case <-ready:
	case lastLine := <-died:
		t.Fatalf("anvil exited before becoming ready, last output: %s", lastLine)
	case <-time.After(120 * time.Second):
		t.Fatal("anvil did not become ready within 120s (slow fork RPC?)")
	}

	url := fmt.Sprintf("http://127.0.0.1:%d", port)
	rpcClient, err := rpc.Dial(url)
	if err != nil {
		t.Fatalf("failed to dial anvil: %v", err)
	}
	ethClient, err := ethclient.Dial(url)
	if err != nil {
		t.Fatalf("failed to dial anvil via ethclient: %v", err)
	}

	return &forkNode{url: url, rpcClient: rpcClient, ethClient: ethClient}
}

func (n *forkNode) setBalance(t *testing.T, address string, wei *big.Int) {
	t.Helper()
	err := n.rpcClient.CallContext(context.Background(), nil, "anvil_setBalance", address, hexutil.EncodeBig(wei))
	if err != nil {
		t.Fatalf("anvil_setBalance failed: %v", err)
	}
}

// setNextBlockTimestamp moves the next block's timestamp to at least the given value,
// clamped above the current head so sequential subtests never rewind the chain
func (n *forkNode) setNextBlockTimestamp(t *testing.T, timestamp int64) {
	t.Helper()
	header, err := n.ethClient.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Fatalf("failed to read latest header: %v", err)
	}
	if head := int64(header.Time); timestamp <= head {
		timestamp = head + 1
	}
	err = n.rpcClient.CallContext(context.Background(), nil, "evm_setNextBlockTimestamp", hexutil.EncodeBig(big.NewInt(timestamp)))
	if err != nil {
		t.Fatalf("evm_setNextBlockTimestamp failed: %v", err)
	}
}

// sendImpersonated sends a transaction from an account anvil unlocks for us (no private key needed)
func (n *forkNode) sendImpersonated(t *testing.T, from, to string, data []byte) {
	t.Helper()
	ctx := context.Background()
	if err := n.rpcClient.CallContext(ctx, nil, "anvil_impersonateAccount", from); err != nil {
		t.Fatalf("anvil_impersonateAccount failed: %v", err)
	}
	defer func() {
		_ = n.rpcClient.CallContext(ctx, nil, "anvil_stopImpersonatingAccount", from)
	}()

	tx := map[string]any{
		"from": from,
		"to":   to,
		"data": hexutil.Encode(data),
	}
	var txHash string
	if err := n.rpcClient.CallContext(ctx, &txHash, "eth_sendTransaction", tx); err != nil {
		t.Fatalf("impersonated eth_sendTransaction failed: %v", err)
	}
	n.requireReceiptStatus(t, txHash, 1)
}

func (n *forkNode) requireReceiptStatus(t *testing.T, txHash string, status uint64) *types.Receipt {
	t.Helper()
	ctx := context.Background()
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		receipt, err := n.ethClient.TransactionReceipt(ctx, geth_common.HexToHash(txHash))
		if err == nil {
			if receipt.Status != status {
				t.Fatalf("tx %s status = %d, want %d", txHash, receipt.Status, status)
			}
			return receipt
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for receipt of %s", txHash)
	return nil
}

// sendTx signs and broadcasts a transaction from the given SDK wallet and requires the given receipt status
func (n *forkNode) sendTx(t *testing.T, wallet *web3_provider.Wallet, txBuilder common.TransactionBuilderFactory, to string, data []byte, value *big.Int) *types.Receipt {
	t.Helper()
	ctx := context.Background()

	toAddress := geth_common.HexToAddress(to)
	builder := txBuilder.New().SetData(data).SetTo(&toAddress)
	if value != nil {
		builder = builder.SetValue(value)
	}
	tx, err := builder.Build(ctx)
	if err != nil {
		t.Fatalf("failed to build tx to %s: %v", to, err)
	}
	signedTx, err := wallet.Sign(tx)
	if err != nil {
		t.Fatalf("failed to sign tx: %v", err)
	}
	if err := wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		t.Fatalf("failed to broadcast tx: %v", err)
	}
	return n.requireReceiptStatus(t, signedTx.Hash().Hex(), 1)
}

// trySendTx is like sendTx but returns the receipt status instead of failing, for negative tests.
// Gas is set manually because eth_estimateGas rejects reverting transactions before broadcast.
func (n *forkNode) trySendTx(t *testing.T, wallet *web3_provider.Wallet, txBuilder common.TransactionBuilderFactory, to string, data []byte) uint64 {
	t.Helper()
	ctx := context.Background()

	toAddress := geth_common.HexToAddress(to)
	tx, err := txBuilder.New().SetData(data).SetTo(&toAddress).SetGas(1_000_000).Build(ctx)
	if err != nil {
		t.Fatalf("failed to build tx to %s: %v", to, err)
	}
	signedTx, err := wallet.Sign(tx)
	if err != nil {
		t.Fatalf("failed to sign tx: %v", err)
	}
	if err := wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		// anvil may reject reverting transactions at broadcast time; any other
		// broadcast failure is an environment problem the test must not hide
		if strings.Contains(strings.ToLower(err.Error()), "revert") {
			t.Logf("broadcast rejected with revert: %v", err)
			return 0
		}
		t.Fatalf("broadcast failed for a reason other than revert: %v", err)
	}
	deadline := time.Now().Add(30 * time.Second)
	for time.Now().Before(deadline) {
		receipt, err := n.ethClient.TransactionReceipt(ctx, signedTx.Hash())
		if err == nil {
			return receipt.Status
		}
		time.Sleep(200 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for receipt of %s", signedTx.Hash().Hex())
	return 0
}

// deployContract deploys the given creation bytecode with ABI-encoded constructor args appended
func (n *forkNode) deployContract(t *testing.T, wallet *web3_provider.Wallet, txBuilder common.TransactionBuilderFactory, creationCode []byte) geth_common.Address {
	t.Helper()
	ctx := context.Background()

	tx, err := txBuilder.New().SetData(creationCode).Build(ctx)
	if err != nil {
		t.Fatalf("failed to build deploy tx: %v", err)
	}
	signedTx, err := wallet.Sign(tx)
	if err != nil {
		t.Fatalf("failed to sign deploy tx: %v", err)
	}
	if err := wallet.BroadcastTransaction(ctx, signedTx); err != nil {
		t.Fatalf("failed to broadcast deploy tx: %v", err)
	}
	receipt := n.requireReceiptStatus(t, signedTx.Hash().Hex(), 1)
	if receipt.ContractAddress == (geth_common.Address{}) {
		t.Fatal("deploy receipt has no contract address")
	}
	return receipt.ContractAddress
}

type contractArtifact struct {
	Abi      json.RawMessage `json:"abi"`
	Bytecode string          `json:"bytecode"`
}

func loadArtifact(t *testing.T, path string) *contractArtifact {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read artifact %s: %v", path, err)
	}
	var artifact contractArtifact
	if err := json.Unmarshal(raw, &artifact); err != nil {
		t.Fatalf("failed to parse artifact %s: %v", path, err)
	}
	return &artifact
}
