package grammar

// nolint: govet
import (
	"github.com/alecthomas/participle/v2"
	"strings"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "TRUE"
	return nil
}

type SimpleIndexDefinition struct {
	IndexName    *string        `@Ident?`
	IndexType    *string        `( "USING" @("BTREE" | "HASH") )?`
	Keys         []*KeyPart     `"(" @@ ( "," @@ )* ")"`
	IndexOptions []*IndexOption `@@*`
}

type PrimaryKeyDefinition struct {
	Constraint   *CheckConstraint `@@?`
	IsPrimary    bool             `"PRIMARY" "KEY"`
	IndexType    *string          `( "USING" @("BTREE" | "HASH") )?`
	Keys         []*KeyPart       `"(" @@ ( "," @@ )* ")"`
	IndexOptions []*IndexOption   `@@*`
}

type UniqueKeyDefinition struct {
	Constraint   *CheckConstraint `@@?`
	IsUnique     bool             `"UNIQUE"`
	IsIndex      bool             `( @"INDEX" | "KEY" )`
	IndexName    *string          `@Ident?`
	IndexType    *string          `( "USING" @("BTREE" | "HASH") )?`
	Keys         []*KeyPart       `"(" @@ ( "," @@ )* ")"`
	IndexOptions []*IndexOption   `@@*`
}

type ForeignKeyDefinition struct {
	Constraint          *CheckConstraint     `@@?`
	IsForeign           bool                 `"FOREIGN" "KEY"`
	IndexName           *string              `@Ident?`
	ColumnNames         []string             `"(" @Ident ( "," @Ident )* ")"`
	ReferenceDefinition *ReferenceDefinition `@@`
}

type SpecialIndexDefinition struct {
	IndexSort    *string        `@( "FULLTEXT" | "SPATIAL" )`
	IsIndex      bool           `( @"INDEX" | "KEY" )`
	IndexName    *string        `@Ident?`
	Keys         []*KeyPart     `"(" @@ ( "," @@ )* ")"`
	IndexOptions []*IndexOption `@@*`
}

type KeyPart struct {
	KeyPartColumn *KeyPartColumn `( @@`
	Expression    *Expression    `| @@ )`
	IsAsc         bool           `( @"ASC" | "DESC" )?`
}

type KeyPartColumn struct {
	Name   string `@Ident`
	Length *int   `( "(" @Number ")" )?`
}

type CheckConstraintDefinition struct {
	Constraint *CheckConstraint `@@?`
	Expression *Expression      ` "CHECK" @@`
	IsEnforced bool             `( @"ENFORCED" | "NOT" "ENFORCED" )?`
}

type CheckConstraint struct {
	Name *string ` "CONSTRAINT" @Ident? `
}

type IndexOption struct {
	KeyBlockSize             *Value  `( @@`
	IndexType                *string ` | "USING" @("BTREE" | "HASH") `
	WithParser               *string ` | "WITH" "PARSER" @Ident`
	Comment                  *string ` | "COMMENT" @String`
	Visible                  bool    ` | ( @"VISIBLE" | "INVISIBLE" )`
	EngineAttribute          *string ` | ( "ENGINE_ATTRIBUTE" "="? @String )`
	SecondaryEngineAttribute *string ` | ( "SECONDARY_ENGINE_ATTRIBUTE" "="? @String ) )`
}

type ColumnDefinition struct {
	Simple   *SimpleColumnDefinition `( @@`
	AsColumn *AsColumnDefinition     `| @@ )`
}

type AsColumnDefinition struct {
	ColumnName                string                     `@Ident`
	DataType                  ColumnDataType             `@@`
	CollationName             *string                    `( "COLLATE" @Ident )?`
	GeneratedAlways           bool                       `@( "GENERATED" "ALWAYS" )?`
	Expression                *Expression                `( "AS" "(" @@ ")" )?`
	IsStored                  bool                       `( @"STORED" `
	IsVirtual                 bool                       `| @"VIRTUAL" )?`
	NotNull                   bool                       `( @( "NOT" "NULL" ) | "NULL" )?`
	Visible                   bool                       `( @"VISIBLE" | "INVISIBLE" )?`
	UniqueKey                 bool                       `@( "UNIQUE" "KEY"? )?`
	PrimaryKey                bool                       `@( "PRIMARY"? "KEY" )?`
	Comment                   *string                    `( "COMMENT" @String )?`
	ReferenceDefinition       *ReferenceDefinition       `@@?`
	CheckConstraintDefinition *CheckConstraintDefinition `@@?`
}

type SimpleColumnDefinition struct {
	ColumnName                string                     `@Ident`
	DataType                  ColumnDataType             `@@`
	NotNull                   bool                       `( @( "NOT" "NULL" ) | "NULL" )?`
	Default                   *ColumnDefault             `( "DEFAULT" @@ )?`
	Visible                   bool                       `( @"VISIBLE" | "INVISIBLE" )?`
	AutoIncrement             bool                       `@"AUTO_INCREMENT"? `
	UniqueKey                 bool                       `@( "UNIQUE" "KEY"? )?`
	PrimaryKey                bool                       `@( "PRIMARY"? "KEY" )?`
	Comment                   *string                    `( "COMMENT" @String )?`
	Collate                   *string                    `( "COLLATE" @Ident )?`
	ColumnFormat              *UppercaseString           `( "COLUMN_FORMAT" @( "FIXED" | "DYNAMIC" | "DEFAULT" ) )?`
	EngineAttribute           *string                    `( "ENGINE_ATTRIBUTE" "="? @String )?`
	SecondaryEngineAttribute  *string                    `( "SECONDARY_ENGINE_ATTRIBUTE" "="? @String )?`
	Storage                   *string                    `( "STORAGE" @( "DISK" | "MEMORY" ) )?`
	ReferenceDefinition       *ReferenceDefinition       `@@?`
	CheckConstraintDefinition *CheckConstraintDefinition `@@?`
}

type ReferenceDefinition struct {
	TableName string           `"REFERENCES" @Ident`
	Keys      []*KeyPart       `"(" @@ ( "," @@ )* ")"`
	Match     *string          `( "MATCH" @( "FULL" | "PARTIAL" | "SIMPLE" ) )?`
	OnDelete  *ReferenceOption `( "ON" "DELETE" @( "RESTRICT" | "CASCADE" | "SET" "NULL" | "NO" "ACTION" | "SET" "DEFAULT" ) )?`
	OnUpdate  *ReferenceOption `( "ON" "UPDATE" @( "RESTRICT" | "CASCADE" | "SET" "NULL" | "NO" "ACTION" | "SET" "DEFAULT" ) )?`
}

type ReferenceOption string

func (r *ReferenceOption) Capture(values []string) error {
	*r = ReferenceOption(strings.Join(values, " "))
	return nil
}

type ColumnDataType struct {
	Bit     *BitDataType     `( @@`
	Integer *IntegerDataType `| @@`
	String  *StringDataType  `| @@`
	EnumSet *EnumDataType    `| @@`
	Blob    *BlobDataType    `| @@`
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

type BlobDataType struct {
	Type *UppercaseString `@("BINARY" | "VARBINARY" | "BLOB" | "TINYBLOB" | "MEDIUMBLOB" | "LONGBLOB")`

	// Precision is mandatory for VARBINARY, and not allowed for TINYBLOB and co, but whatever
	Precision *int `( "(" @Number ")" )?`
}

type EnumDataType struct {
	IsSet         bool     `( "ENUM" | @"SET" )`
	Values        []string `( "(" @String ( "," @String )* ")" )?`
	CharacterSet  *string  `( "CHARACTER" "SET" @Ident )?`
	CollationName *string  `( "COLLATE" @Ident )?`
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

var (
	parser = participle.MustBuild(
		&CreateTable{},
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
		participle.CaseInsensitive("Keyword"),
		participle.Elide("Comment"),
		// Need to solve left recursion detection first, if possible.
		participle.UseLookahead(1),
	)
)

func Parse(s string) (*CreateTable, error) {
	sql := &CreateTable{}
	err := parser.ParseString("", s, sql)
	return sql, err
}
