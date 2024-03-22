package types_test

import (
	"testing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"

	"github.com/skip-mev/slinky/x/mm2/types"
)

func TestValidateBasic(t *testing.T) {
	type testCase struct {
		name      string
		params    types.Params
		expectErr bool
	}

	testCases := []testCase{
		{
			name:      "valid default params",
			params:    types.DefaultParams(),
			expectErr: false,
		},
		{
			name: "valid multiple authorities",
			params: types.Params{
				MarketAuthorities: []string{authtypes.NewModuleAddress(authtypes.ModuleName).String(), types.DefaultMarketAuthority},
				Version:           10,
			},
			expectErr: false,
		},
		{
			name: "invalid duplicate authority",
			params: types.Params{
				MarketAuthorities: []string{types.DefaultMarketAuthority, types.DefaultMarketAuthority},
				Version:           10,
			},
			expectErr: true,
		},
		{
			name: "invalid authority string",
			params: types.Params{
				MarketAuthorities: []string{"incorrect"},
				Version:           10,
			},
			expectErr: true,
		},
		{
			name: "invalid nil authority",
			params: types.Params{
				MarketAuthorities: nil,
				Version:           10,
			},
			expectErr: true,
		},
		{
			name:      "invalid empty params",
			params:    types.Params{},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.params.ValidateBasic()
			if tc.expectErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
