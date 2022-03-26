package mysql

import "github.com/spf13/cobra"

var MysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql related commands",
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "parse and dump mysql schema",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	MysqlCmd.AddCommand(schemaCmd)
}
