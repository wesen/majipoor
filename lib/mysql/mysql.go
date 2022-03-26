package mysql

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"reflect"
	"time"
)

type MysqlSlaveStatus struct {
	RetrievedGtidSet string `db:"Retrieved_Gtid_Set"`
}

type MysqlGlobalVariables struct {
	GtidMode       string `mysql:"gtid_mode"`
	LogBin         string `mysql:"log_bin"`
	BinlogFormat   string `mysql:"binlog_format"`
	BinlogRowImage string `mysql:"binlog_row_image"`
	ServerUuid     string `mysql:"server_uuid"`
}

type MysqlDB struct {
	db *sqlx.DB
}

func (md *MysqlDB) Close() error {
	return md.db.Close()
}

func NewMysqlDB(connectionString string) (*MysqlDB, error) {
	db, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return &MysqlDB{db: db}, nil
}

func getMysqlVariable(db *sqlx.DB, variableName string) (string, error) {
	res := db.QueryRow(fmt.Sprintf("SHOW GLOBAL VARIABLES LIKE '%s'", variableName))
	var val string
	var variable string
	err := res.Scan(&variable, &val)
	return val, err
}

func (md *MysqlDB) GetMysqlGlobalVariables() (*MysqlGlobalVariables, error) {
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
			value, err := getMysqlVariable(md.db, v)
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

	return &config, nil
}

func (md *MysqlDB) GetMysqlSlaveStatus() (*MysqlSlaveStatus, error) {
	slaveStatus := &MysqlSlaveStatus{}
	err := md.db.Unsafe().Get(slaveStatus, "SHOW SLAVE status")
	if err == sql.ErrNoRows {
		err = nil
	}
	return slaveStatus, err
}
