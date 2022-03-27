// Package mysql contains commands to interact with the mysql server
//
// Inspired by pg_chameleon, Copyright (c) 2016-2020 Federico Campoli
package mysql

import (
	"github.com/huandu/go-sqlbuilder"
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

func getMysqlConnectionString(username string, password string) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		username, password,
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		viper.GetString("mysql.database"))
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "parse and dump mysql schema",
	Run: func(cmd *cobra.Command, args []string) {
		connectionString := getMysqlConnectionString(
			viper.GetString("mysql.username"),
			viper.GetString("mysql.password"))
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
		tables, err := db.GetTables(schema, viper.GetStringSlice("mysql.limit-tables"),
			viper.GetStringSlice("mysql.skip-tables"))
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
				log.Info().Str("table", table).Str("column", c.ColumnName).Msg("Found column")
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
		connectionString := getMysqlConnectionString(
			viper.GetString("mysql.username"),
			viper.GetString("mysql.password"))
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

var createReplicaUserCmd = &cobra.Command{
	Use:   "create-replica-user",
	Short: "Create a replica user",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

		rootUsername, _ := cmd.Flags().GetString("mysql-root-username")
		rootPassword, _ := cmd.Flags().GetString("mysql-root-password")
		connectionString := getMysqlConnectionString(rootUsername, rootPassword)
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

		sb := sqlbuilder.Select("user").From("mysql.user")
		replicaUsername := viper.GetString("mysql.username")
		sb.Where(sb.Equal("user", replicaUsername))
		sb2 := sqlbuilder.Buildf("SELECT EXISTS(%v)", sb)
		sql_, args_ := sb2.Build()

		if force {
			sql_ = fmt.Sprintf("DROP USER IF EXISTS %s", replicaUsername)
			if dryRun {
				log.Info().Str("sql", sql_).Msg("Force deletion of user")
			} else {
				_, err = db.Exec(sql_)
				if err != nil {
					log.Fatal().Err(err).Msg("Could not delete user")
				}
				log.Info().Str("username", replicaUsername).Msg("Deleting user")
			}
		} else {
			if dryRun {
				log.Info().Str("sql", sql_).Interface("args", args_).Msg("Checking if user exists")
			} else {
				var exists bool
				err = db.Db.QueryRow(sql_, args_...).Scan(&exists)
				if err != nil {
					log.Fatal().Err(err).Msg("Could not check if user exists")
				}
				if exists {
					log.Fatal().Str("username", replicaUsername).Msg("User already exists")
				}
			}
		}

		sql_ = fmt.Sprintf("CREATE USER %s", replicaUsername)
		if dryRun {
			log.Info().Str("sql", sql_).Msg("Creating user")
		} else {
			_, err = db.Exec(sql_)
			if err != nil {
				log.Fatal().Err(err).Msg("Could not create user")
			}
			log.Info().Str("username", replicaUsername).Msg("Creating user")
		}

		sql_ = fmt.Sprintf("SET PASSWORD FOR %s=PASSWORD('%s')",
			replicaUsername, viper.GetString("mysql.password"))
		if dryRun {
			log.Info().Str("sql", sql_).Msg("Setting password")
		} else {
			_, err = db.Exec(sql_)
			if err != nil {
				log.Fatal().Err(err).Msg("Could not set password")
			}
			log.Info().Str("username", replicaUsername).Msg("Setting password")
		}

		replicaSchema := viper.GetString("mysql.schema")
		sql_ = fmt.Sprintf("GRANT ALL ON %s.* TO '%s'",
			replicaSchema, replicaUsername)

		if dryRun {
			log.Info().Str("sql", sql_).Msg("Granting privileges")
		} else {
			_, err = db.Exec(sql_)

			if err != nil {
				log.Fatal().Err(err).Msg("Could not grant schema privileges")
			}
			log.Info().Str("username", replicaUsername).Str("schema", replicaSchema).Msg("Replica schema privileges granted")
		}

		for _, v := range []string{
			"GRANT RELOAD ON *.* TO '%s'",
			"GRANT REPLICATION CLIENT ON *.* TO '%s'",
			"GRANT REPLICATION SLAVE ON *.* TO '%s'",
		} {
			sql_ = fmt.Sprintf(v, replicaUsername)
			if dryRun {
				log.Info().Str("sql", sql_).Msg("Granting replication privileges")
			} else {
				_, err = db.Exec(sql_)
				if err != nil {
					log.Fatal().Err(err).Msg("Could not grant replication privileges")
				}
				log.Info().Str("username", replicaUsername).Msg("Replication privileges granted")
			}
		}

		sql_ = "FLUSH PRIVILEGES"
		if dryRun {
			log.Info().Str("sql", sql_).Msg("Flushing privileges")
		} else {
			_, err = db.Exec(sql_)

			if err != nil {
				log.Fatal().Err(err).Msg("Could not flush privileges")
			}
			log.Info().Str("username", replicaUsername).Msg("Privileges flushed")
		}
	},
}

// TODO(manuel) On postgresql side
// CREATE USER usr_replica WITH PASSWORD 'replica';
// CREATE DATABASE db_replica WITH OWNER usr_replica;

// TODO(manuel) Add command to check my.cnf file
//binlog_format= ROW
//binlog_row_image=FULL
//log-bin = mysql-bin
//server-id = 1
//expire_logs_days = 10

func init() {
	createReplicaUserCmd.Flags().Bool("dry-run", false, "Dry run")
	createReplicaUserCmd.Flags().Bool("force", false, "Force recreation")

	MysqlCmd.AddCommand(schemaCmd, checkConfigCmd, createReplicaUserCmd)

}
