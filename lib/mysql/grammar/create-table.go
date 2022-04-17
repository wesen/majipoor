package grammar

type TableName string

type CreateTable struct {
	Temporary        bool                     `"CREATE" ( @"TEMPORARY" )?  "TABLE"`
	Name             string                   `@( Ident ( "." Ident )* )`
	IfNotExists      bool                     `@( "IF" "NOT" "EXISTS" )?`
	CreateDefinition []*CreateTableDefinition `"(" ( @@ ( "," @@ )* )? ")"`
	TableOptions     []TableOption            `@@*`
	PartitionOptions *PartitionOptions        ` ( "PARTITION" "BY" @@ )?`
}

type CreateTableDefinition struct {
	ColumnDefinition          *ColumnDefinition          `( @@ `
	SimpleIndexDefinition     *SimpleIndexDefinition     ` | ( "INDEX" | "KEY" ) @@`
	SpecialIndexDefinition    *SpecialIndexDefinition    ` | @@`
	PrimaryKeyDefinition      *PrimaryKeyDefinition      ` | @@`
	UniqueKeyDefinition       *UniqueKeyDefinition       ` | @@`
	ForeignKeyDefinition      *ForeignKeyDefinition      ` | @@`
	CheckConstraintDefinition *CheckConstraintDefinition ` | @@ )`
}

type TableOption struct {
	AutoExtendSize           *int        `( "AUTOEXTEND_SIZE" "="? @Number`
	AutoIncrement            *int        ` | "AUTO_INCREMENT" "="? @Number`
	AvgRowLength             *int        ` | "AVG_ROW_LENGTH" "="? @Number`
	CharacterSet             *string     ` | "DEFAULT"? "CHARACTER" "SET" "="? @Ident`
	Checksum                 *int        ` | "CHECKSUM" "="? @Number`
	Collation                *string     ` | "DEFAULT"? "COLLATE" "="? @Ident`
	Comment                  *string     ` | "COMMENT" "="? @String`
	Compression              *string     ` | "COMPRESSION" "="? @String`
	Connection               *string     ` | "CONNECTION" "="? @String`
	DataDirectory            *string     ` | "DATA" "DIRECTORY" "="? @String`
	IndexDirectory           *string     ` | "INDEX" "DIRECTORY" "="? @String`
	DelayKeyWrite            *int        ` | "DELAY_KEY_WRITE" "="? @Number`
	Encryption               *string     ` | "ENCRYPTION" "="? @String`
	Engine                   *string     ` | "ENGINE" "="? @Ident`
	EngineAttribute          *string     ` | "ENGINE_ATTRIBUTE" "="? @String `
	InsertMethod             *string     ` | "INSERT_METHOD" "="? @("NO" | "FIRST" | "LAST")`
	SecondaryEngineAttribute *string     ` | "SECONDARY_ENGINE_ATTRIBUTE" "="? @String`
	KeyBlockSize             *int        ` | "KEY_BLOCK_SIZE" "="? @Number`
	MaxRows                  *int        ` | "MAX_ROWS" "="? @Number`
	MinRows                  *int        ` | "MIN_ROWS" "="? @Number`
	PackKeys                 *string     ` | "PACK_KEYS" "="? @( Number | "DEFAULT" )`
	Password                 *string     ` | "PASSWORD" "="? @String`
	RowFormat                *string     ` | "ROW_FORMAT" "="? @( "DEFAULT" | "DYNAMIC" | "FIXED" | "COMPRESSED" | "REDUNDANT" | "COMPACT" )`
	StatsAutoRecalc          *string     ` | "STATS_AUTO_RECALC" "="? @( "DEFAULT" | Number )`
	StatsPersistent          *string     ` | "STATS_PERSISTENT" "="? @( "DEFAULT" | Number )`
	StatsSamplePages         *int        ` | "STATS_SAMPLE_PAGES" "="? @Number`
	TableSpace               *TableSpace ` | @@`
	Union                    []string    ` | "UNION" "="? "(" @Ident ( "," @Ident )* ")" )`
}

type TableSpace struct {
	Name            string `"TABLESPACE" @Ident`
	IsDiskStorage   bool   `( "STORAGE" @"DISK" `
	IsMemoryStorage bool   `| "STORAGE" @"MEMORY" )?`
}
