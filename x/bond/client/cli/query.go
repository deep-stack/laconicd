package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	bondtypes "git.vdb.to/cerc-io/laconicd/x/bond"
)

// GetQueryBondList implements the bond lists query command.
func GetQueryBondList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List bonds.",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get bond list .

Example:
$ %s query %s list
`,
				version.AppName, bondtypes.ModuleName,
			),
		),
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}
			queryClient := bondtypes.NewQueryClient(clientCtx)
			res, err := queryClient.Bonds(cmd.Context(), &bondtypes.QueryBondsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
