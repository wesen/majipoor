package gen

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"majipoor/lib/helpers"
	"majipoor/lib/mysql"
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
		log.Debug().Msgf("Generating fake data for table %s", name)
		t := schema_gen.GenerateTable(name)
		createStatement := t.TableDefinition()
		fmt.Printf(createStatement)

		apply, _ := cmd.Flags().GetBool("apply")
		if apply {
			connectionString := helpers.GetReplicaMysqlConnectionString()
			log.Debug().Str("mysql-connection-string", connectionString).Msg("Connecting to mysql")
			db, err := mysql.NewMysqlDB(connectionString)
			if err != nil {
				log.Fatal().Err(err).Msg("Could not connect to database")
			}

			defer func() {
				err := db.Close()
				if err != nil {
					log.Error().Err(err).Msg("Could not close database connection")
				}
			}()

			force, _ := cmd.Flags().GetBool("force")
			if force {
				log.Info().Str("table", name).Msg("Dropping table")
				_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", name))
				if err != nil {
					log.Fatal().Err(err).Str("table", name).Msg("Could not drop table")
				}
			}

			_, err = db.Exec(createStatement)
			if err != nil {
				log.Fatal().Err(err).Str("table", name).Msg("Could not create table")
			}
		}
	},
}

func init() {
	rand.Seed(time.Now().Unix())
	GenCmd.AddCommand(genTableCmd)
	genTableCmd.Flags().Bool("apply", false, "Apply generated data to database")
	genTableCmd.Flags().Bool("force", false, "Force overwrite of existing table")
	genTableCmd.Flags().String("table", "", "Table name")
}
