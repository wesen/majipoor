package mysql

import (
	"github.com/alecthomas/repr"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"majipoor/lib/mysql/grammar"
)

var parseCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parse a MySQL statement",
	Run: func(cmd *cobra.Command, args []string) {
		for _, v := range args {
			sql, err := grammar.Parse(v)
			if err != nil {
				log.Warn().Err(err).Str("sql", v).Msg("failed to parse")
				continue
			}
			repr.Println(sql, repr.Indent("  "), repr.OmitEmpty(true))
		}
	},
}

func init() {
	MysqlCmd.AddCommand(parseCmd)
}
