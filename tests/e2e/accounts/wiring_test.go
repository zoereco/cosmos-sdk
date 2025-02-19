//go:build app_v1

package accounts

import (
	"testing"

	"cosmossdk.io/core/header"
	counterv1 "cosmossdk.io/x/accounts/testing/counter/v1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

// TestDependencies aims to test wiring between different account components,
// inherited from the runtime, specifically:
// - address codec
// - binary codec
// - header service
// - gas service
func TestDependencies(t *testing.T) {
	app := setupApp(t)
	ak := app.AccountsKeeper
	ctx := sdk.NewContext(app.CommitMultiStore(), false, app.Logger()).WithHeaderInfo(header.Info{ChainID: "chain-id"})

	_, counterAddr, err := ak.Init(ctx, "counter", accCreator, &counterv1.MsgInit{
		InitialValue: 0,
	})
	require.NoError(t, err)
	// test dependencies
	r, err := ak.Execute(ctx, counterAddr, []byte("test"), &counterv1.MsgTestDependencies{})
	require.NoError(t, err)
	res := r.(*counterv1.MsgTestDependenciesResponse)

	// test gas
	require.NotZero(t, res.BeforeGas)
	require.NotZero(t, res.AfterGas)
	require.Equal(t, uint64(10), res.AfterGas-res.BeforeGas)

	// test header service
	require.Equal(t, ctx.HeaderInfo().ChainID, res.ChainId)

	// test address codec
	wantAddr, err := app.AuthKeeper.AddressCodec().BytesToString(counterAddr)
	require.NoError(t, err)
	require.Equal(t, wantAddr, res.Address)
}
