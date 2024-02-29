package gql

import "github.com/spf13/cobra"

// AddGQLFlags adds gql flags for
func AddGQLFlags(cmd *cobra.Command) *cobra.Command {
	// Add flags for GQL server.
	cmd.PersistentFlags().Bool("gql-server", false, "Start GQL server.")
	cmd.PersistentFlags().Bool("gql-playground", false, "Enable GQL playground.")
	cmd.PersistentFlags().String("gql-playground-api-base", "", "GQL API base path to use in GQL playground.")
	cmd.PersistentFlags().String("gql-port", "9473", "Port to use for the GQL server.")
	cmd.PersistentFlags().String("log-file", "", "File to tail for GQL 'getLogs' API.")

	return cmd
}
