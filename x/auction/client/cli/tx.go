package cli

import (
	"encoding/hex"
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	wnsUtils "git.vdb.to/cerc-io/laconic2d/utils"
	auctiontypes "git.vdb.to/cerc-io/laconic2d/x/auction"
)

// GetTxCmd returns transaction commands for this module.
func GetTxCmd() *cobra.Command {
	auctionTxCmd := &cobra.Command{
		Use:                        auctiontypes.ModuleName,
		Short:                      "Auction transactions subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	auctionTxCmd.AddCommand(
		GetCmdCommitBid(),
		GetCmdRevealBid(),
	)

	return auctionTxCmd
}

// GetCmdCommitBid is the CLI command for committing a bid.
func GetCmdCommitBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "commit-bid [auction-id] [bid-amount]",
		Short: "Commit sealed bid",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			bidAmount, err := sdk.ParseCoinNormalized(args[1])
			if err != nil {
				return err
			}

			mnemonic, err := wnsUtils.GenerateMnemonic()
			if err != nil {
				return err
			}

			chainId := viper.GetString("chain-id")
			auctionId := args[0]

			reveal := map[string]interface{}{
				"chainId":       chainId,
				"auctionId":     auctionId,
				"bidderAddress": clientCtx.GetFromAddress().String(),
				"bidAmount":     bidAmount.String(),
				"noise":         mnemonic,
			}

			commitHash, content, err := wnsUtils.GenerateHash(reveal)
			if err != nil {
				return err
			}

			// Save reveal file.
			err = os.WriteFile(fmt.Sprintf("%s-%s.json", clientCtx.GetFromName(), commitHash), content, 0o600)
			if err != nil {
				return err
			}

			msg := auctiontypes.NewMsgCommitBid(auctionId, commitHash, clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

// GetCmdRevealBid is the CLI command for revealing a bid.
func GetCmdRevealBid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reveal-bid [auction-id] [reveal-file-path]",
		Short: "Reveal bid",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			auctionId := args[0]
			revealFilePath := args[1]

			revealBytes, err := os.ReadFile(revealFilePath)
			if err != nil {
				return err
			}

			msg := auctiontypes.NewMsgRevealBid(auctionId, hex.EncodeToString(revealBytes), clientCtx.GetFromAddress())
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
