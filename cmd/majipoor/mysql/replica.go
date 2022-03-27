package mysql

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"majipoor/lib/helpers"
	"majipoor/lib/mysql"
)

var createReplicaUserCmd = &cobra.Command{
	Use:   "create-replica-user",
	Short: "Create a replica user",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

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

		err = db.CreateReplicaUser(mysql.CreateReplicaUserSettings{
			Force:    force,
			DryRun:   dryRun,
			Schema:   viper.GetString("mysql.database"),
			Username: viper.GetString("mysql.username"),
			Password: viper.GetString("mysql.password"),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create replica user")
		}
	},
}

var createReplicaDatabaseCmd = &cobra.Command{
	Use:   "create-replica-database",
	Short: "Create a replica database",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		force, _ := cmd.Flags().GetBool("force")

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

		err = db.CreateReplicaDatabase(mysql.CreateReplicaDatabaseSettings{
			Force:    force,
			DryRun:   dryRun,
			Database: viper.GetString("mysql.database"),
		})
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create replica database")
		}
	},
}
