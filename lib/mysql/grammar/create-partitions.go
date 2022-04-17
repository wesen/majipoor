package grammar

type PartitionOptions struct {
	HashPartition        *HashPartition         `( @@ `
	KeyPartition         *KeyPartition          `| @@`
	RangePartition       *RangePartition        `| "RANGE" @@ `
	ListPartition        *ListPartition         `| "LIST" @@ )`
	Partitions           *int                   `( "PARTITIONS" @Number )?`
	SubpartitionByHash   *HashPartition         `( "SUBPARTITION" "BY" @@`
	SubPartitionByKey    *KeyPartition          `| "SUBPARTITION" "BY" @@ )?`
	SubPartitions        *int                   `( "SUBPARTITIONS" @Number )?`
	PartitionDefinitions []*PartitionDefinition `( "(" @@ ( "," @@ )* ")" )?`
}

type HashPartition struct {
	IsLinear   bool       `@"LINEAR"?`
	Expression Expression `"HASH" "(" @@ ")"`
}

type KeyPartition struct {
	IsLinear  bool     `@"LINEAR"? "KEY"`
	Algorithm *int     `("ALGORITHM" "=" @Number)?`
	Columns   []string `"(" @Ident ( "," @Ident )* ")"`
}

type RangePartition struct {
	Expression *Expression `( "(" @@ ")"`
	Columns    []string    `| "COLUMNS" "(" @Ident ( "," @Ident )* ")" )`
}

type ListPartition struct {
	Expression *Expression `( "(" @@ ")"`
	Columns    []string    `| "COLUMNS" "(" @Ident ( "," @Ident )* ")" )`
}

type PartitionDefinition struct {
	Name           string          `"PARTITION" @Ident`
	ValuesLessThan *ValuesLessThan `( "VALUES" "LESS" "THAN" @@`
	ValuesIn       []Value         ` | "VALUES" "IN" "(" @@ ( "," @@ )* ")" )?`
}

type ValuesLessThan struct {
	IsMaxValue bool        `( @"MAXVALUE"?`
	Expression *Expression `| "(" @@ ")"`
	Values     []Value     `| "(" @@ ( "," @@ )* ")" )`
}
