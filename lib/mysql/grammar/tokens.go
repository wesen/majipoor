package grammar

import (
	"fmt"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

var keywords = []string{
	"ACTION", "ALGORITHM", "ALL", "ALWAYS", "AND", "ANY", "AS", "AUTOEXTEND_SIZE",
	"AUTO_INCREMENT", "AVG_ROW_LENGTH", "BETWEEN", "BTREE", "BY", "CASCADE", "CHARACTER",
	"CHECK", "CHECKSUM", "COLLATE", "COLUMNS", "COLUMN_FORMAT", "COMMENT", "COMPACT",
	"COMPRESSED", "COMPRESSION", "CONNECTION", "CONSTRAINT", "CREATE", "DATA", "DEFAULT",
	"DELAY_KEY_WRITE", "DELETE", "DIRECTORY", "DISK", "DISTINCT", "DYNAMIC", "ENCRYPTION",
	"ENFORCED", "ENGINE", "ENGINE_ATTRIBUTE", "EXCEPT", "EXISTS", "FALSE", "FIRST", "FIXED",
	"FOREIGN", "FROM", "FULLTEXT", "GENERATED", "GROUP", "HASH", "HAVING", "IF", "IN", "INDEX",
	"INSERT_METHOD", "INTERSECT", "INVISIBLE", "IS", "KEY", "KEY_BLOCK_SIZE", "LAST", "LIKE",
	"LIMIT", "LINEAR", "LIST", "MATCH", "MAXVALUE", "MAX_ROWS", "MEMORY", "MINUS", "MIN_ROWS",
	"NO", "NOT", "NULL", "OFFSET", "ON", "OR", "ORDER", "PACK_KEYS", "PARSER", "PARTIAL", "PARTITION",
	"PRIMARY", "RANGE", "REDUNDANT", "REFERENCES", "RESTRICT", "ROW_FORMAT", "SECONDARY_ENGINE_ATTRIBUTE",
	"SELECT", "SET", "SIMPLE", "SOME", "SPATIAL", "STATS_AUTO_RECALC", "STATS_PERSISTENT", "STATS_SAMPLE_PAGES",
	"STORAGE", "STORED", "SUBPARTITION", "SUBPARTITIONS", "TABLE", "TABLESPACE", "TEMPORARY",
	"TOP", "TRUE", "TYPE", "UNION", "UNIQUE", "UPDATE", "USING", "VIEW", "VIRTUAL", "VISIBLE", "WHERE", "WITH",

	"TRUE", "FALSE",
	"CURRENT_TIMESTAMP", "LOCALTIME", "NOW", "LOCALTIMESTAMP",

	"NATIONAL",

	"INSTANT", "INPLACE", "COPY", "ALGORITHM", "CHANGE", "AFTER", "FIRST", "DROP", "CONVERT", "DISABLE", "ENABLE",
	"DISCARD". "IMPORT", "LOCK", "RENAME", "MODIFY", "SHARED", "EXCLUSIVE" ,"WITHOUT", "VALIDATION", "TO",
	"TRUNCATE", "DISCARD", "COALESCE", "REORGANIZE", "ANALYZE", "OPTIMIZE", "REBUILD", "REPAIR", "REMOVE",
}

var types = []string{
	"CHAR", "BYTE", "BINARY", "NCHAR", "VARCHAR", "VARBINARY", "ENUM", "SET", "TEXT",
	"DATE", "TIME", "TIMESTAMP", "DATETIME", "YEAR",
	"FLOAT", "DOUBLE", "REAL", "PRECISION",
	"BIT",
	"BOOL", "BOOLEAN",
	"BLOB", "TINYBLOB", "TINYTEXT", "MEDIUMBLOB", "MEDIUMTEXT", "LONGBLOB", "LONGTEXT",
	"INTEGER", "INT", "SMALLINT", "TINYINT", "MEDIUMINT", "BIGINT",
	"UNSIGNED", "ZEROFILL",
	"NUMERIC", "DECIMAL", "DEC", "FIXED",
	"GEOMETRY", "POINT", "LINESTRING", "POLYGON", "MULTIPOINT", "MULTILINESTRING", "MULTIPOLYGON", "GEOMETRYCOLLECTION",
}

var sqlLexer = lexer.MustSimple([]lexer.Rule{
	{Name: "Comment", Pattern: ` //.*|/\*.*?\*/`},
	{
		Name:    `Keyword`,
		Pattern: fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(append(keywords, types...), "|")),
	},
	{
		Name: "whitespace", Pattern: `\s+`,
	},
	{
		Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`,
	},
	{
		Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`,
	},
	{
		Name: `String`, Pattern: `'[^']*'|"[^"]*"`,
	},
	{
		Name: `Operators`, Pattern: `<>|!=|<=|>=|[-+*/%,.()=<>]`,
	},
},
)
