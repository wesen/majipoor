package grammar

import (
	"github.com/stretchr/testify/assert"
	require "github.com/stretchr/testify/require"
	"testing"
)

func TestParseCaseInsensitive(t *testing.T) {
	for _, s := range []string{
		"CREATE TABLE foobar ( id INT )",
		"create table foobar ( id INT )",
	} {
		ast, err := Parse(s)
		if assert.Nil(t, err) && assert.NotNil(t, ast) {
			assert.Equal(t, ast.Name, "foobar")
		}
	}
}

func TestParseCaseComment(t *testing.T) {
	for _, s := range []string{
		"CREATE TABLE foobar ( id INT ) /* generated by server */",
	} {
		ast, err := Parse(s)
		if assert.Nil(t, err) && assert.NotNil(t, ast) {
			assert.Equal(t, "foobar", ast.Name)
		}
	}
}

func TestParseBit(t *testing.T) {
	ast, err := Parse("CREATE TABLE foobar ( bitColumn BIT )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, "bitColumn", definition.ColumnName)
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Nil(t, definition.ColumnDefinition.DataType.Bit.Precision)
		require.Nil(t, definition.ColumnDefinition.DataType.Integer)
	}
	ast, err = Parse("CREATE TABLE foobar ( bitColumn BIT(5) )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, "bitColumn", definition.ColumnName)
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Equal(t, 5, *definition.ColumnDefinition.DataType.Bit.Precision)
		require.Nil(t, definition.ColumnDefinition.DataType.Integer)
	}
	ast, err = Parse("CREATE TABLE foobar ( bitColumn BIT (5 ) )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, "bitColumn", definition.ColumnName)
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Equal(t, 5, *definition.ColumnDefinition.DataType.Bit.Precision)
		require.Nil(t, definition.ColumnDefinition.DataType.Integer)
	}
}

func TestParseInteger(t *testing.T) {
	for _, s := range []string{
		"CREATE TABLE foobar ( intColumn TINYINT )",
		"CREATE TABLE foobar ( intColumn tinyint )",
	} {
		ast, err := Parse(s)
		if assert.Nil(t, err) && assert.NotNil(t, ast) {
			definition := ast.CreateDefinition[0]
			assert.Equal(t, "intColumn", definition.ColumnName)
			require.NotNil(t, definition.ColumnDefinition.DataType.Integer)
			require.Nil(t, definition.ColumnDefinition.DataType.Integer.Precision)
			require.NotNil(t, definition.ColumnDefinition.DataType.Integer.Type)
			require.Equal(t, UppercaseString("TINYINT"), *definition.ColumnDefinition.DataType.Integer.Type)
			require.Equal(t, false, definition.ColumnDefinition.DataType.Integer.Unsigned)
			require.Nil(t, definition.ColumnDefinition.DataType.Bit)
		}
	}

	ast, err := Parse("CREATE TABLE foobar ( intColumn TINYINT(4) UNSIGNED ZEROFILL )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		integer := ast.CreateDefinition[0].ColumnDefinition.DataType.Integer
		require.Equal(t, UppercaseString("TINYINT"), *integer.Type)
		require.True(t, integer.Unsigned)
		require.True(t, integer.Zerofill)
	}
}

func TestParseBool(t *testing.T) {
	for _, s := range []string{
		"CREATE TABLE foobar ( boolColumn BOOL )",
		"CREATE TABLE foobar ( boolColumn boolEAN )",
		"CREATE TABLE foobar ( boolColumn BOOLEAN )",
	} {
		ast, err := Parse(s)
		if assert.Nil(t, err) && assert.NotNil(t, ast) {
			definition := ast.CreateDefinition[0]
			assert.Equal(t, "boolColumn", definition.ColumnName)
			require.True(t, definition.ColumnDefinition.DataType.Bool)
		}
	}
}

func TestParseMultipleColumns(t *testing.T) {
	ast, err := Parse(`CREATE TABLE foobar ( 
	bitColumn BIT(5) , 
    foobar VARCHAR(30) NOT NULL DEFAULT 'foo' ,
     blop INT DEFAULT 2
    )
`)
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		require.Equal(t, 3, len(ast.CreateDefinition))
		definition := ast.CreateDefinition[0]
		require.False(t, definition.ColumnDefinition.NotNull)
		require.Equal(t, "bitColumn", definition.ColumnName)
		require.Nil(t, definition.ColumnDefinition.Default)
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)

		definition = ast.CreateDefinition[1]
		require.True(t, definition.ColumnDefinition.NotNull)
		require.Equal(t, "foobar", definition.ColumnName)
		require.NotNil(t, definition.ColumnDefinition.Default)
		require.Equal(t, "foo", *definition.ColumnDefinition.Default.String)

		definition = ast.CreateDefinition[2]
		require.False(t, definition.ColumnDefinition.NotNull)
		require.Equal(t, "blop", definition.ColumnName)
		require.NotNil(t, definition.ColumnDefinition.Default)
		require.Equal(t, 2.0, *definition.ColumnDefinition.Default.Number)
		require.NotNil(t, definition.ColumnDefinition.DataType.Integer)
		require.Equal(t, UppercaseString("INT"), *definition.ColumnDefinition.DataType.Integer.Type)
	}
}

func TestParseColumnOptions(t *testing.T) {
	ast, err := Parse("CREATE TABLE foobar ( bitColumn BIT(5) NOT NULL DEFAULT 1 VISIBLE AUTO_INCREMENT UNIQUE KEY PRIMARY KEY COMMENT 'comment')")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0].ColumnDefinition
		require.True(t, definition.NotNull)
		require.NotNil(t, definition.Default)
		require.True(t, definition.Visible)
		require.True(t, definition.AutoIncrement)
		require.True(t, definition.UniqueKey)
		require.True(t, definition.PrimaryKey)
		require.Equal(t, *definition.Comment, "comment")
	}

	ast, err = Parse(`CREATE TABLE foobar ( 
	bitColumn BIT(5) 
      NOT NULL
      DEFAULT 1
      INVISIBLE
      AUTO_INCREMENT UNIQUE KEY
      PRIMARY KEY
      COMMENT 'comment and comment'
      COLLATE utf8_bin
      COLUMN_FORMAT DEfaULT
    )
`)
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0].ColumnDefinition
		require.True(t, definition.NotNull)
		require.NotNil(t, definition.Default)
		require.False(t, definition.Visible)
		require.True(t, definition.AutoIncrement)
		require.True(t, definition.UniqueKey)
		require.True(t, definition.PrimaryKey)
		require.Equal(t, *definition.Comment, "comment and comment")
		require.Equal(t, *definition.Collate, "utf8_bin")
		require.Equal(t, *definition.ColumnFormat, UppercaseString("DEFAULT"))
	}
	ast, err = Parse(`CREATE TABLE foobar ( 
	bitColumn BIT(5) 
      COLUMN_FORMAT blabl
    )
`)
	assert.NotNil(t, err)
}