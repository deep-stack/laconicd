package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	bondtypes "git.vdb.to/cerc-io/laconicd/x/bond"
)

// NewCreateBondCmd is the CLI command for creating a bond.
// Used in e2e tests
func NewCreateBondCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [amount]",
		Short: "Create bond.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			coin, err := sdk.ParseCoinNormalized(args[0])
			if err != nil {
				return err
			}

			msg := bondtypes.NewMsgCreateBond(sdk.NewCoins(coin), clientCtx.GetFromAddress())
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
