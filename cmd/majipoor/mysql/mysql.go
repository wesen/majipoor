// Package mysql contains commands to interact with the mysql server
//
// Inspired by pg_chameleon, Copyright (c) 2016-2020 Federico Campoli
package mysql

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
	"reflect"
	"time"
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

/**
  def __check_mysql_config(self):
      """
          The method check if the mysql configuration is compatible with the replica requirements.
          If all the configuration requirements are met then the return value is True.
          Otherwise is false.
          The parameters checked are
          log_bin - ON if the binary log is enabled
          binlog_format - must be ROW , otherwise the replica won't get the data
          binlog_row_image - must be FULL, otherwise the row image will be incomplete

      """
      if self.gtid_enable:
          sql_log_bin = """SHOW GLOBAL VARIABLES LIKE 'gtid_mode';"""
          self.cursor_buffered.execute(sql_log_bin)
          variable_check = self.cursor_buffered.fetchone()
          if variable_check:
              gtid_mode = variable_check["Value"]
              if gtid_mode.upper() == 'ON':
                  self.gtid_mode = True
                  sql_uuid = """SHOW SLAVE STATUS;"""
                  self.cursor_buffered.execute(sql_uuid)
                  slave_status = self.cursor_buffered.fetchall()
                  if len(slave_status)>0:
                      gtid_set=slave_status[0]["Retrieved_Gtid_Set"]
                  else:
                      sql_uuid = """SHOW GLOBAL VARIABLES LIKE 'server_uuid';"""
                      self.cursor_buffered.execute(sql_uuid)
                      server_uuid = self.cursor_buffered.fetchone()
                      gtid_set = server_uuid["Value"]
                  self.gtid_uuid = gtid_set.split(':')[0]

          else:
              self.gtid_mode = False
      else:
          self.gtid_mode = False

      sql_log_bin = """SHOW GLOBAL VARIABLES LIKE 'log_bin';"""
      self.cursor_buffered.execute(sql_log_bin)
      variable_check = self.cursor_buffered.fetchone()
      log_bin = variable_check["Value"]

      sql_log_bin = """SHOW GLOBAL VARIABLES LIKE 'binlog_format';"""
      self.cursor_buffered.execute(sql_log_bin)
      variable_check = self.cursor_buffered.fetchone()
      binlog_format = variable_check["Value"]

      sql_log_bin = """SHOW GLOBAL VARIABLES LIKE 'binlog_row_image';"""
      self.cursor_buffered.execute(sql_log_bin)
      variable_check = self.cursor_buffered.fetchone()
      if variable_check:
          binlog_row_image = variable_check["Value"]
      else:
          binlog_row_image = 'FULL'

      if log_bin.upper() == 'ON' and binlog_format.upper() == 'ROW' and binlog_row_image.upper() == 'FULL':
          self.replica_possible = True
      else:
          self.replica_possible = False
          self.pg_engine.set_source_status("error")
          self.logger.error("The MySQL configuration does not allow the replica. Exiting now")
          self.logger.error("Source settings - log_bin %s, binlog_format %s, binlog_row_image %s" % (log_bin.upper(),  binlog_format.upper(), binlog_row_image.upper() ))
          self.logger.error("Mandatory settings - log_bin ON, binlog_format ROW, binlog_row_image FULL (only for MySQL 5.6+) ")
          sys.exit()

*/

type MysqlSlaveStatus struct {
	RetrievedGtidSet string `db:"Retrieved_Gtid_Set"`
}

type MysqlGlobalVariables struct {
	GtidMode       string `mysql:"gtid_mode"`
	LogBin         string `mysql:"log_bin"`
	BinlogFormat   string `mysql:"binlog_format"`
	BinlogRowImage string `mysql:"binlog_row_image"`
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "parse and dump mysql schema",
	Run: func(cmd *cobra.Command, args []string) {
		connectionString := getMysqlConnectionString()
		log.Debug().Str("mysql-connection-string", connectionString).Msg("Connecting to mysql")
		db, err := sqlx.Open("mysql", connectionString)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connect to database")
		}

		defer func() {
			err := db.Close()
			if err != nil {
				log.Error().Err(err).Msg("Could not close database connection")
			}
		}()
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxIdleConns(10)
		db.SetMaxOpenConns(10)

		var config MysqlGlobalVariables
		t := reflect.TypeOf(config)
		ps := reflect.ValueOf(&config)
		s := ps.Elem()
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			sf := s.Field(i)

			v, ok := f.Tag.Lookup("mysql")
			if ok {
				_log := log.With().Str("field", f.Name).Str("mysql-variable", v).Logger()
				value, err := getMysqlVariable(db, v)
				if err != nil && err != sql.ErrNoRows {
					_log.Error().Err(err).Msg("Could not query variable")
					continue
				}
				if sf.CanSet() {
					if sf.Kind() == reflect.String {
						sf.SetString(value)
					} else {
						_log.Error().Str("value", value).Msg("Field is not a string")
					}
				}
				_log.Info().Str("value", value).Send()
			} else {
				log.Warn().Str("field", f.Name).Msg("No mysql tag found")
			}
		}

		log.Info().Interface("config", config).Send()

		slaveStatus := &MysqlSlaveStatus{}
		err = db.Unsafe().Get(slaveStatus, "SHOW SLAVE status")
		if err == sql.ErrNoRows {
			log.Info().Msg("No slave status found")
		} else if err != nil {
			log.Fatal().Err(err).Msg("Could not get slave status")
		} else {
			log.Info().Interface("slave_status", slaveStatus).Send()
		}

	},
}

func getMysqlVariable(db *sqlx.DB, variableName string) (string, error) {
	res := db.QueryRow(fmt.Sprintf("SHOW GLOBAL VARIABLES LIKE '%s'", variableName))
	var val string
	var variable string
	err := res.Scan(&variable, &val)
	return val, err
}

func init() {
	MysqlCmd.AddCommand(schemaCmd)
}
