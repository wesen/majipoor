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
		name, _ := cmd.Flags().GetString("name")
		t := schema_gen.GenerateTable(name)
		fmt.Printf(t.TableDefinition())

		apply, _ := cmd.Flags().GetBool("apply")
		if apply {

		}
	},
}

func init() {
	rand.Seed(time.Now().Unix())
	GenCmd.AddCommand(genTableCmd)
	genTableCmd.Flags().Bool("apply", false, "Apply generated data to database")
	genTableCmd.Flags().String("table", "", "Table name")
}
