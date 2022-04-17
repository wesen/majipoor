package grammar

type AlterTable struct {
	Name             string                 `"ALTER" "TABLE" `
	AlterOptions     []AlterOption          ` @@ ( "," @@ )*`
	PartitionOptions []AlterPartitionOption `@@ ( "," @@ )*`
}

type AlterOption struct {
	AddSingleColumn              *AddSingleColumn           `( "ADD" "COLUMN"? @@`
	AddMultipleColumns           []AddColumnDefinition      `| "ADD" "COLUMN"? "(" @@ ( "," @@ )* ")"`
	AddSimpleIndex               *SimpleIndexDefinition     `| "ADD" ( "INDEX" | "KEY" ) @@`
	AddPrimaryKeyDefinition      *PrimaryKeyDefinition      `| "ADD" @@`
	AddUniqueKeyDefinition       *UniqueKeyDefinition       `| "ADD" @@`
	AddForeignKeyDefinition      *ForeignKeyDefinition      `| "ADD" @@`
	AddSpecialIndexDefinition    *SpecialIndexDefinition    `| "ADD" @@`
	AddCheckConstraintDefinition *CheckConstraintDefinition `| "ADD" @@`

	DropCheckConstraint  *string               `| "DROP" ( "CHECK" | "CONSTRAINT") @Ident`
	AlterCheckConstraint *AlterCheckConstraint `| "ALTER" ( "CHECK" | "CONSTRAINT" ) @@`
	Algorithm            *string               `| "ALGORITHM" "="? @( "DEFAULT" | "INSTANT" | "INPLACE" | "COPY" )`
}

type AddSingleColumn struct {
	Name             string           `@Ident`
	ColumnDefinition ColumnDefinition `@@`
	IsFirst          bool             `@"FIRST"?`
	After            *string          `( "AFTER" @Ident )?`
}

type AddColumnDefinition struct {
	Name             string           `@Ident`
	ColumnDefinition ColumnDefinition `@@`
}

type AlterPartitionOption struct {
}

type AlterCheckConstraint struct {
	Name       string `@Ident`
	IsEnforced bool   `( @"ENFORCED" | "NOT" "ENFORCED" )`
}
