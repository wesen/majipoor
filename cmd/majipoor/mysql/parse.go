package mysql

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"majipoor/lib/mysql/grammar"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a MySQL statement",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			_, _ = grammar.Parse(v)

		}
		log.Info().Strs("args", args).Send()
	},
}

func init() {
	MysqlCmd.AddCommand(parseCmd)
}
