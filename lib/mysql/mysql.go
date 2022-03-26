package mysql

import (
	"database/sql"
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"reflect"
	"strings"
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

func (md *MysqlDB) GetTables(schema string, limitTables []string, skipTables []string) ([]string, error) {
	var tables []string
	sb := sqlbuilder.Select("TABLE_NAME").From("information_schema.TABLES")
	sb.Where(sb.Equal("TABLE_TYPE", "BASE TABLE"))
	sb.Where(sb.Equal("TABLE_SCHEMA", schema))
	sql_, args := sb.Build()

	checkLimitTables := false
	limitTableCheck := map[string]bool{}
	skipTableCheck := map[string]bool{}

	if len(limitTables) > 0 {
		checkLimitTables = true
		for _, v := range limitTables {
			limitTableCheck[v] = true
		}
	}
	for _, v := range skipTables {
		skipTableCheck[v] = true
	}

	rows, err := md.db.Query(sql_, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var name string
		err = rows.Scan(&name)
		if err != nil {
			return nil, err
		}
		if _, ok := skipTableCheck[name]; ok {
			continue
		}

		if checkLimitTables {
			if _, ok := limitTableCheck[name]; !ok {
				continue
			}
		}
		tables = append(tables, name)
	}

	return tables, nil
}

type ColumnMetadata struct {
	ColumnName             string  `db:"column_name"`
	ColumnDefault          *string `db:"column_default"`
	OrdinalPosition        int     `db:"ordinal_position"`
	DataType               string  `db:"data_type"`
	ColumnType             string  `db:"column_type"`
	CharacterMaximumLength *int    `db:"character_maximum_length"`
	Extra                  string  `db:"extra"`
	ColumnKey              string  `db:"column_key"`
	IsNullable             string  `db:"is_nullable"`
	NumericPrecision       *int    `db:"numeric_precision"`
	NumericScale           *int    `db:"numeric_scale"`
	EnumList               *string `db:"enum_list"`
}

func (md *MysqlDB) GetTableMetadata(schema string, table string) ([]*ColumnMetadata, error) {
	var metadatas []*ColumnMetadata
	sb := sqlbuilder.Select("column_name", "column_default", "ordinal_position",
		"data_type", "column_type", "character_maximum_length", "extra", "column_key",
		"is_nullable", "numeric_precision", "numeric_scale",
		`CASE
	         WHEN data_type="enum"
	     THEN
	         SUBSTRING(COLUMN_TYPE,5)
	     END AS enum_list`).
		From("information_schema.COLUMNS")
	sb.Where(sb.Equal("TABLE_SCHEMA", schema))
	sb.Where(sb.Equal("TABLE_NAME", table))
	sb.OrderBy("ordinal_position")

	sql_, args := sb.Build()

	rows, err := md.db.Queryx(sql_, args...)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()
	for rows.Next() {
		var tableMetadata ColumnMetadata
		err = rows.StructScan(&tableMetadata)
		if err != nil {
			return nil, err
		}
		metadatas = append(metadatas, &tableMetadata)
	}

	return metadatas, nil
}

var spatialDatatypes = []string{
	"point", "geometry", "linestring", "polygon", "multipoint", "multilinestring", "geometrycollection",
}

var hexTypes = []string{
	"blob", "tinyblob", "mediumblob", "longblob", "binary", "varbinary",
}

// contains checks if needle is in haystack. This is horrible from a O() standpoint,
// but we don't care
func contains(needle string, haystack []string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}

var defaultCharacterSet = "utf8"

func (c *ColumnMetadata) getSelectCSVStatement() string {
	if contains(c.DataType, hexTypes) {
		return fmt.Sprintf("hex(%s)", c.ColumnName)
	}
	if c.DataType == "bit" {
		return fmt.Sprintf("cast(%s AS unsigned)", c.ColumnName)
	}
	if contains(c.DataType, []string{"datetime", "timestamp", "date"}) {
		return fmt.Sprintf("nullif(%s, cast(\"0000-00-00 00:00:00\" AS date))", c.ColumnName)
	}
	if contains(c.DataType, spatialDatatypes) {
		return fmt.Sprintf("ST_AsText(%s)", c.ColumnName)
	}

	return fmt.Sprintf("cast(%s AS char CHARACTER SET %s)", c.ColumnName, defaultCharacterSet)
}

func mapValues(m map[string]string) []string {
	res := []string{}
	for _, v := range m {
		res = append(res, v)
	}
	return res
}

func GetSelectCSVSatement(table string, columns []*ColumnMetadata) string {
	_ = table
	selects := map[string]string{}
	selectCsvs := map[string]string{}
	for _, c := range columns {
		statement := c.getSelectCSVStatement()
		selects[c.ColumnName] = statement
		selectCsvs[c.ColumnName] = fmt.Sprintf("COALESCE(REPLACE(%s, '\"', '\"\"'), 'NULL')", statement)
		log.Debug().Str("column", c.ColumnName).Str("type", c.ColumnType).Str("select", selects[c.ColumnName]).Send()
	}
	return fmt.Sprintf("REPLACE(CONCAT('\"',CONCAT_WS('\",\"',%s),'\"'),'\"NULL\"','NULL')",
		strings.Join(mapValues(selectCsvs), ",\n"))
}
