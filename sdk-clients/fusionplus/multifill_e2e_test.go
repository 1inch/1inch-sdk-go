package fusionplus

import (
	"context"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/1inch/1inch-sdk-go/v5/common"
	"github.com/1inch/1inch-sdk-go/v5/constants"
)

// captureExecutor records the requests the SDK makes so a test can assert the
// exact wire bytes without a live relayer.
type captureExecutor struct {
	Payloads []common.RequestPayload
}

func (e *captureExecutor) ExecuteRequest(_ context.Context, payload common.RequestPayload, _ any) error {
	e.Payloads = append(e.Payloads, payload)
	return nil
}

// TestPendingSecretIndexes covers the secret-selection logic that the maker
// monitor loop uses. The maker must reveal the secret each ready fill names in
// its Idx field, once per index, regardless of the order fills arrive in.
func TestPendingSecretIndexes(t *testing.T) {
	tests := []struct {
		name        string
		fillIdxs    []float32
		revealed    map[int]bool
		secretCount int
		want        []int
	}{
		{
			name:        "single resolver fills whole order needs the last secret",
			fillIdxs:    []float32{3},
			secretCount: 4,
			want:        []int{3},
		},
		{
			name:        "fills arrive out of order",
			fillIdxs:    []float32{2, 0},
			secretCount: 4,
			want:        []int{0, 2},
		},
		{
			name:        "duplicates in one batch collapse",
			fillIdxs:    []float32{1, 1, 2},
			secretCount: 4,
			want:        []int{1, 2},
		},
		{
			name:        "already-revealed indexes are skipped",
			fillIdxs:    []float32{0, 1, 2},
			revealed:    map[int]bool{0: true, 1: true},
			secretCount: 4,
			want:        []int{2},
		},
		{
			name:        "out-of-range indexes are dropped",
			fillIdxs:    []float32{4, -1, 3},
			secretCount: 4,
			want:        []int{3},
		},
		{
			name:        "nothing ready yields nothing",
			fillIdxs:    nil,
			secretCount: 4,
			want:        nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fills := make([]ReadyToAcceptSecretFill, len(tc.fillIdxs))
			for i, idx := range tc.fillIdxs {
				fills[i] = ReadyToAcceptSecretFill{Idx: idx}
			}
			got := PendingSecretIndexes(fills, tc.revealed, tc.secretCount)
			assert.Equal(t, tc.want, got)
		})
	}
}

// multiFillOrderParams builds order params for an n-secret multiple-fill order,
// with a real hashlock and matching secret hashes.
func multiFillOrderParams(t *testing.T, n int) OrderParams {
	t.Helper()
	secrets := genSecrets(n)
	secretHashes := make([]string, n)
	for i, s := range secrets {
		h, err := HashSecret(s)
		require.NoError(t, err)
		secretHashes[i] = h
	}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	hashLock, err := ForMultipleFills(leaves)
	require.NoError(t, err)

	return OrderParams{
		HashLock:     hashLock,
		SecretHashes: secretHashes,
		Preset:       Fast,
		Receiver:     constants.ZeroAddress,
		Nonce:        big.NewInt(1),
	}
}

func newTestClient(e *captureExecutor) *Client {
	return &Client{
		api:    api{httpExecutor: e},
		Wallet: &stubWallet{addr: gethcommon.HexToAddress("0x4444444444444444444444444444444444444444")},
	}
}

// TestPlaceOrder_MultiFill_SendsSecretHashes is the wire-level end-to-end check
// for Finding 1. A multiple-fill order must now be accepted by PlaceOrder and
// the submitted body must carry the secret hashes the relayer needs to rebuild
// the Merkle tree. Before the fix, PlaceOrder rejected the order outright and
// never transmitted the hashes.
func TestPlaceOrder_MultiFill_SendsSecretHashes(t *testing.T) {
	orderParams := multiFillOrderParams(t, 4)
	executor := &captureExecutor{}
	client := newTestClient(executor)

	orderHash, err := client.PlaceOrder(
		context.Background(),
		takerAssetTestQuoteParams(collisionToken),
		takerAssetTestQuote(),
		orderParams,
		client.Wallet,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, orderHash)

	require.Len(t, executor.Payloads, 1)
	body := string(executor.Payloads[0].Body)
	assert.Contains(t, executor.Payloads[0].Path, "/submit")
	assert.Contains(t, body, `"secretHashes"`)
	for _, h := range orderParams.SecretHashes {
		assert.Contains(t, body, h, "submitted body must carry every secret hash")
	}
}

// TestPlaceOrder_SingleFill_OmitsSecretHashes confirms a single-fill order still
// omits the secret hashes (its hashlock is the lone secret hash).
func TestPlaceOrder_SingleFill_OmitsSecretHashes(t *testing.T) {
	secret := genSecrets(1)[0]
	hashLock, err := ForSingleFill(secret)
	require.NoError(t, err)
	hash, err := HashSecret(secret)
	require.NoError(t, err)

	orderParams := OrderParams{
		HashLock:     hashLock,
		SecretHashes: []string{hash},
		Preset:       Fast,
		Receiver:     constants.ZeroAddress,
		Nonce:        big.NewInt(1),
	}
	executor := &captureExecutor{}
	client := newTestClient(executor)

	_, err = client.PlaceOrder(
		context.Background(),
		takerAssetTestQuoteParams(collisionToken),
		takerAssetTestQuote(),
		orderParams,
		client.Wallet,
	)
	require.NoError(t, err)
	require.Len(t, executor.Payloads, 1)
	assert.NotContains(t, string(executor.Payloads[0].Body), `"secretHashes"`)
}

// TestPlaceOrder_MultiFill_Guards covers the two rejection paths that mirror the
// JS SDK contract: a preset that disallows multiple fills cannot carry more than
// one secret hash, and a multiple-fill order must carry exactly parts+1 hashes.
func TestPlaceOrder_MultiFill_Guards(t *testing.T) {
	t.Run("preset disallows multiple fills", func(t *testing.T) {
		quote := takerAssetTestQuote()
		preset := quote.Presets.Fast
		preset.AllowMultipleFills = false
		quote.Presets.Fast = preset

		executor := &captureExecutor{}
		client := newTestClient(executor)
		_, err := client.PlaceOrder(
			context.Background(),
			takerAssetTestQuoteParams(collisionToken),
			quote,
			multiFillOrderParams(t, 4),
			client.Wallet,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "allows multiple fills")
		assert.Empty(t, executor.Payloads, "no order submitted when the guard rejects")
	})

	t.Run("secret hash count must match the hashlock parts", func(t *testing.T) {
		orderParams := multiFillOrderParams(t, 4) // hashlock parts = 3, so 4 hashes expected
		orderParams.SecretHashes = orderParams.SecretHashes[:3]

		executor := &captureExecutor{}
		client := newTestClient(executor)
		_, err := client.PlaceOrder(
			context.Background(),
			takerAssetTestQuoteParams(collisionToken),
			takerAssetTestQuote(),
			orderParams,
			client.Wallet,
		)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "parts count plus one")
		assert.Empty(t, executor.Payloads)
	})
}
