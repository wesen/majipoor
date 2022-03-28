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

func TestParseBit(t *testing.T) {
	ast, err := Parse("CREATE TABLE foobar ( bitColumn BIT )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, definition.ColumnName, "bitColumn")
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Nil(t, definition.ColumnDefinition.DataType.Bit.Precision)
		require.Nil(t, definition.ColumnDefinition.DataType.Integer)
	}
	ast, err = Parse("CREATE TABLE foobar ( bitColumn BIT(5) )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, definition.ColumnName, "bitColumn")
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Equal(t, *definition.ColumnDefinition.DataType.Bit.Precision, 5)
		require.Nil(t, definition.ColumnDefinition.DataType.Integer)
	}
	ast, err = Parse("CREATE TABLE foobar ( bitColumn BIT (5 ) )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		definition := ast.CreateDefinition[0]
		assert.Equal(t, definition.ColumnName, "bitColumn")
		require.NotNil(t, definition.ColumnDefinition.DataType.Bit)
		require.Equal(t, *definition.ColumnDefinition.DataType.Bit.Precision, 5)
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
			assert.Equal(t, definition.ColumnName, "intColumn")
			require.NotNil(t, definition.ColumnDefinition.DataType.Integer)
			require.Nil(t, definition.ColumnDefinition.DataType.Integer.Precision)
			require.NotNil(t, definition.ColumnDefinition.DataType.Integer.Type)
			require.Equal(t, *definition.ColumnDefinition.DataType.Integer.Type, SqlType("TINYINT"))
			require.Equal(t, definition.ColumnDefinition.DataType.Integer.Unsigned, false)
			require.Nil(t, definition.ColumnDefinition.DataType.Bit)
		}
	}

	ast, err := Parse("CREATE TABLE foobar ( intColumn TINYINT(4) UNSIGNED ZEROFILL )")
	if assert.Nil(t, err) && assert.NotNil(t, ast) {
		integer := ast.CreateDefinition[0].ColumnDefinition.DataType.Integer
		require.Equal(t, *integer.Type, SqlType("TINYINT"))
		require.True(t, integer.Unsigned)
		require.True(t, integer.Zerofill)
	}
}
