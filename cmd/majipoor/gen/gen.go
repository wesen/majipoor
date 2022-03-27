package gen

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	schema_gen "majipoor/lib/schema-gen"
)

var GenCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate fake SQL data",
}

var genTableCmd = &cobra.Command{
	Use:   "table",
	Short: "Generate fake SQL data",
	Run: func(cmd *cobra.Command, args []string) {
		t := schema_gen.GenerateTable()
		log.Info().Interface("table", t).Msg("Generated table")
	},
}

func init() {
	GenCmd.AddCommand(genTableCmd)
}
