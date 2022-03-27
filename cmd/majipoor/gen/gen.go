package gen

import (
	"fmt"
	"github.com/spf13/cobra"
	schema_gen "majipoor/lib/schema-gen"
	"math/rand"
	"time"
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
		fmt.Printf(t.TableDefinition())
	},
}

func init() {
	rand.Seed(time.Now().Unix())
	GenCmd.AddCommand(genTableCmd)
}
