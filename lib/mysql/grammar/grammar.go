package grammar

// nolint: govet
import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"strings"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type Create struct {
	Temporary        bool                `"CREATE" ( @"TEMPORARY" )?  "TABLE"`
	Name             string              `@Ident`
	IfNotExists      bool                `@( "IF" "NOT" "EXISTS" )?`
	CreateDefinition []*CreateDefinition `"(" ( @@ ( "," @@ )* )? ")"`
	//TableOptions     *TableOptions       `@@`
	//PartitionOptions *PartitionOptions   `@@`
}

type CreateDefinition struct {
	ColumnName       string            `@Ident`
	ColumnDefinition *ColumnDefinition `@@`
}

type ColumnDefinition struct {
	DataType      ColumnDataType   `@@`
	NotNull       bool             `( @( "NOT" "NULL" ) | "NULL" )?`
	Default       *ColumnDefault   `( "DEFAULT" @@ )?`
	Visible       bool             `( @"VISIBLE" | "INVISIBLE" )?`
	AutoIncrement bool             `@"AUTO_INCREMENT"? `
	UniqueKey     bool             `@( "UNIQUE" "KEY"? )?`
	PrimaryKey    bool             `@( "PRIMARY"? "KEY" )?`
	Comment       *string          `( "COMMENT" @String )?`
	Collate       *string          `( "COLLATE" @Ident )?`
	ColumnFormat  *UppercaseString `( "COLUMN_FORMAT" @( "FIXED" | "DYNAMIC" | "DEFAULT" ) )?`
}

type ColumnDataType struct {
	Bit     *BitDataType     `( @@`
	Integer *IntegerDataType `| @@`
	String  *StringDataType  `| @@`
	Bool    bool             `| @( "BOOL" | "BOOLEAN" )`
	Last    bool             `)`
}

type BitDataType struct {
	Precision *int `"BIT" ( "(" @Number ")" )?`
}

type IntegerDataType struct {
	Type      *UppercaseString `@("TINYINT" | "SMALLINT" | "MEDIUMINT" | "INT" | "BIGINT" | "INTEGER")`
	Precision *int             `( "(" @Number ")" )?`
	Unsigned  bool             `@"UNSIGNED"?`
	Zerofill  bool             `@"ZEROFILL"?`
}

type StringDataType struct {
	IsNational bool             `@"NATIONAL"?`
	Type       *UppercaseString `@("CHAR" | "VARCHAR" | "TEXT" | "TINYTEXT" | "MEDIUMTEXT" | "LONGTEXT")`

	// Precision is mandatory for VARCHAR, and not allowed for TINYTEXT and co, but whatever
	Precision     *int    `( "(" @Number ")" )?`
	CharacterSet  *string `( "CHARACTER" "SET" @Ident )?`
	CollationName *string `( "COLLATE" @Ident )?`
}

type UppercaseString string

func (s *UppercaseString) Capture(values []string) error {
	*s = UppercaseString(strings.ToUpper(values[0]))
	return nil
}

type ColumnDefault struct {
	Number  *float64 `( @Number`
	String  *string  ` | @String`
	Boolean *Boolean ` | @("TRUE" | "FALSE")`
	Null    bool     ` | @"NULL"`
	Array   *Array   ` | @@ )`
}

type TableOptions struct {
}

type PartitionOptions struct {
}

// Select based on http://www.h2database.com/html/grammar.html
type Select struct {
	Top        *Term             `"SELECT" ( "TOP" @@ )?`
	Distinct   bool              `(  @"DISTINCT"`
	All        bool              ` | @"ALL" )?`
	Expression *SelectExpression `@@`
	From       *From             `"FROM" @@`
	Limit      *Expression       `( "LIMIT" @@ )?`
	Offset     *Expression       `( "OFFSET" @@ )?`
	GroupBy    *Expression       `( "GROUP" "BY" @@ )?`
}

type From struct {
	TableExpressions []*TableExpression `@@ ( "," @@ )*`
	Where            *Expression        `( "WHERE" @@ )?`
}

type TableExpression struct {
	Table  string        `( @Ident ( "." @Ident )*`
	Select *Select       `  | "(" @@ ")"`
	Values []*Expression `  | "VALUES" "(" @@ ( "," @@ )* ")")`
	As     string        `( "AS" @Ident )?`
}

type SelectExpression struct {
	All         bool                 `  @"*"`
	Expressions []*AliasedExpression `| @@ ( "," @@ )*`
}

type AliasedExpression struct {
	Expression *Expression `@@`
	As         string      `( "AS" @Ident )?`
}

type Expression struct {
	Or []*OrCondition `@@ ( "OR" @@ )*`
}

type OrCondition struct {
	And []*Condition `@@ ( "AND" @@ )*`
}

type Condition struct {
	Operand *ConditionOperand `  @@`
	Not     *Condition        `| "NOT" @@`
	Exists  *Select           `| "EXISTS" "(" @@ ")"`
}

type ConditionOperand struct {
	Operand      *Operand      `@@`
	ConditionRHS *ConditionRHS `@@?`
}

type ConditionRHS struct {
	Compare *Compare `  @@`
	Is      *Is      `| "IS" @@`
	Between *Between `| "BETWEEN" @@`
	In      *In      `| "IN" "(" @@ ")"`
	Like    *Like    `| "LIKE" @@`
}

type Compare struct {
	Operator string         `@( "<>" | "<=" | ">=" | "=" | "<" | ">" | "!=" )`
	Operand  *Operand       `(  @@`
	Select   *CompareSelect ` | @@ )`
}

type CompareSelect struct {
	All    bool    `(  @"ALL"`
	Any    bool    ` | @"ANY"`
	Some   bool    ` | @"SOME" )`
	Select *Select `"(" @@ ")"`
}

type Like struct {
	Not     bool     `[ @"NOT" ]`
	Operand *Operand `@@`
}

type Is struct {
	Not          bool     `[ @"NOT" ]`
	Null         bool     `( @"NULL"`
	DistinctFrom *Operand `  | "DISTINCT" "FROM" @@ )`
}

type Between struct {
	Start *Operand `@@`
	End   *Operand `"AND" @@`
}

type In struct {
	Select      *Select       `  @@`
	Expressions []*Expression `| @@ ( "," @@ )*`
}

type Operand struct {
	Summand []*Summand `@@ ( "|" "|" @@ )*`
}

type Summand struct {
	LHS *Factor `@@`
	Op  string  `[ @("+" | "-")`
	RHS *Factor `  @@ ]`
}

type Factor struct {
	LHS *Term  `@@`
	Op  string `( @("*" | "/" | "%")`
	RHS *Term  `  @@ )?`
}

type Term struct {
	Select        *Select     `  @@`
	Value         *Value      `| @@`
	SymbolRef     *SymbolRef  `| @@`
	SubExpression *Expression `| "(" @@ ")"`
}

type SymbolRef struct {
	Symbol     string        `@Ident @( "." Ident )*`
	Parameters []*Expression `( "(" @@ ( "," @@ )* ")" )?`
}

type Value struct {
	Wildcard bool     `(  @"*"`
	Number   *float64 ` | @Number`
	String   *string  ` | @String`
	Boolean  *Boolean ` | @("TRUE" | "FALSE")`
	Null     bool     ` | @"NULL"`
	Array    *Array   ` | @@ )`
}

type Array struct {
	Expressions []*Expression `"(" @@ ( "," @@ )* ")"`
}

var (
	keywords = []string{
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
	}
	types = []string{
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
	sqlLexer = lexer.MustSimple([]lexer.Rule{
		{Name: "Comment", Pattern: `//.*|/\*.*?\*/`},
		{
			Name:    `Keyword`,
			Pattern: fmt.Sprintf(`(?i)\b(%s)\b`, strings.Join(append(keywords, types...), "|")),
		},
		{Name: "whitespace", Pattern: `\s+`},
		{Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`},
		{Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`},
		{Name: `String`, Pattern: `'[^']*'|"[^"]*"`},
		{Name: `Operators`, Pattern: `<>|!=|<=|>=|[-+*/%,.()=<>]`},
	})
	parser = participle.MustBuild(
		&Create{},
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.Elide("Comment"),
		// Need to solve left recursion detection first, if possible.
		participle.UseLookahead(1),
	)
)

func Parse(s string) (*Create, error) {
	sql := &Create{}
	err := parser.ParseString("", s, sql)
	return sql, err
}
