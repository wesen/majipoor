package grammar

type TableExpression struct {
	Table  string        `@( Ident ( "." Ident )*`
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
