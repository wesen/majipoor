package mysql

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"majipoor/lib/helpers"
	"majipoor/lib/mysql"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "parse and dump mysql schema",
	Run: func(cmd *cobra.Command, args []string) {
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

		database := viper.GetString("mysql.database")
		tables, err := db.GetTables(database, viper.GetStringSlice("mysql.limit-tables"),
			viper.GetStringSlice("mysql.skip-tables"))
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get tables")
		}
		for _, table := range tables {
			columns, err := db.GetTableMetadata(database, table)
			if err != nil {
				log.Fatal().Err(err).Str("table", table).Msg("Could not get table metadata")
			}
			log.Info().Str("table", table).Msg("Found table")
			for _, c := range columns {
				log.Info().Str("table", table).Str("column", c.ColumnName).Msg("Found column")
			}

			stmt := mysql.GetSelectCSVSatement(table, columns)
			fmt.Println(stmt)
		}
	},
}
