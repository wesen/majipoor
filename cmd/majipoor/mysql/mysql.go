// Package mysql contains commands to interact with the mysql server
//
// Inspired by pg_chameleon, Copyright (c) 2016-2020 Federico Campoli
package mysql

import (
	"github.com/spf13/cobra"
	"majipoor/lib/mysql"
)

import (
	"fmt"
	_ "github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var MysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql related commands",
}

func getMysqlConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "parse and dump mysql schema",
	Run: func(cmd *cobra.Command, args []string) {
		connectionString := getMysqlConnectionString()
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

		schema := viper.GetString("mysql.schema")
		tables, err := db.GetTables(schema)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get tables")
		}
		for _, table := range tables {
			columns, err := db.GetTableMetadata(schema, table)
			if err != nil {
				log.Fatal().Err(err).Str("table", table).Msg("Could not get table metadata")
			}
			log.Info().Str("table", table).Msg("Found table")
			for _, c := range columns {
				log.Info().Str("table", table).Str("column", c.ColumnName).Msg("Found table")
			}

			stmt := mysql.GetSelectCSVSatement(table, columns)
			fmt.Println(stmt)
		}
	},
}

var checkConfigCmd = &cobra.Command{
	Use:   "check-config",
	Short: "parse and dump mysql config",
	Run: func(cmd *cobra.Command, args []string) {
		connectionString := getMysqlConnectionString()
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

		config, err := db.GetMysqlGlobalVariables()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not get mysql global variables")
		}

		if config.LogBin == "ON" && config.BinlogFormat == "ROW" && config.BinlogRowImage == "FULL" {
			log.Info().Msg("Replica possible")
		} else {
			log.Error().Msg("Replica not possible")
		}
	},
}

func init() {
	MysqlCmd.AddCommand(schemaCmd, checkConfigCmd)
}
