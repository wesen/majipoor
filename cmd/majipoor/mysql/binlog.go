package mysql

import (
	"context"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var binlogCmd = &cobra.Command{
	Use:   "binlog",
	Short: "Subscribe to mysql binlog",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := replication.BinlogSyncerConfig{
			ServerID: 100,
			Flavor:   "mysql",
			Host:     viper.GetString("mysql.host"),
			Port:     uint16(viper.GetInt("mysql.port")),
			User:     viper.GetString("mysql.username"),
			Password: viper.GetString("mysql.password"),
		}

		syncer := replication.NewBinlogSyncer(cfg)
		gtidSet, err := mysql.ParseGTIDSet("mysql", "")
		if err != nil {
			log.Fatal().Err(err).Msg("Could not parse gtid set")
		}
		streamer, err := syncer.StartSyncGTID(gtidSet)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not start binlog sync")
		}

		for {
			ev, err := streamer.GetEvent(context.Background())
			if err != nil {
				log.Warn().Err(err).Send()
			}
			// Dump event
			ev.Dump(os.Stdout)
		}
	},
}

func init() {
	MysqlCmd.AddCommand(binlogCmd)
}
