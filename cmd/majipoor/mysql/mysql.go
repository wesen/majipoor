// Package mysql contains commands to interact with the mysql server
//
// Inspired by pg_chameleon, Copyright (c) 2016-2020 Federico Campoli
package mysql

import (
	"database/sql"
	"github.com/spf13/cobra"
	"majipoor/lib/helpers"
	"majipoor/lib/mysql"
)

import (
	"fmt"
	_ "github.com/go-mysql-org/go-mysql/replication"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
)

var MysqlCmd = &cobra.Command{
	Use:   "mysql",
	Short: "mysql related commands",
}

var checkConfigCmd = &cobra.Command{
	Use:   "check-config",
	Short: "parse and dump mysql config",
	Run: func(cmd *cobra.Command, args []string) {
		connectionString := helpers.GetRootMysqlConnectionString()
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

		log.Info().Interface("mysql-config", config).Send()

		slaveStatus, err := db.GetMysqlSlaveStatus()
		if err == sql.ErrNoRows {
			log.Warn().Msg("No slave status found")
		} else if err != nil {
			log.Fatal().Err(err).Msg("Could not get slave status")
		} else {
			log.Info().Interface("slave-status", slaveStatus).Send()
		}

		if config.LogBin == "ON" && config.BinlogFormat == "ROW" && config.BinlogRowImage == "FULL" && config.GtidMode == "ON" {
			log.Info().Msg("Replica possible")
		} else {
			log.Error().Msg("Replica not possible")
			fmt.Println(`Please add the following entries to your my.cnf file (or binlog.cnf under /etc/mysql/mysql.conf.d/
to enable binary logging:

[mysqld]
binlog_format= ROW
binlog_row_image=FULL
log-bin = mysql-bin
server-id = 1
expire_logs_days = 10
gtid_mode=ON
enforce-gtid-consistency=ON`)
		}
	},
}

// TODO(manuel) On postgresql side
// CREATE USER usr_replica WITH PASSWORD 'replica';
// CREATE DATABASE db_replica WITH OWNER usr_replica;

// TODO(manuel) Build a tool to generate full schemas and test data, so that we can test replication against a real setup
// - this should generate schemas, fake data for the schemas, inserts, updates, deletes
// - it should also generate DDL statements (alter, drop, etc...)
//
// This could output SQL files, not just execute it against a DB itself.
// Now, if we want to generate further data in the future, for more interactive testing,
// should we store the generated schema in a config file?
// Should we be able to generate test data for an existing schema?

// TODO(manuel) Run the schema dump against the ttc database and see which types we get

// TODO(manuel) Gather which indexes to create when inspecting the schema

// TODO create a test framework using a docker test DB to test binlog streaming

func init() {
	createReplicaUserCmd.Flags().Bool("dry-run", false, "Dry run")
	createReplicaUserCmd.Flags().Bool("force", false, "Force recreation")

	createReplicaDatabaseCmd.Flags().Bool("dry-run", false, "Dry run")
	createReplicaDatabaseCmd.Flags().Bool("force", false, "Force recreation")

	MysqlCmd.AddCommand(schemaCmd, checkConfigCmd, createReplicaUserCmd, createReplicaDatabaseCmd)

}
