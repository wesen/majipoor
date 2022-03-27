package schema_gen

import (
	"fmt"
	"math/rand"
)

type Index struct {
	Name        string
	ColumnNames []string
	IsUnique    bool
}

type Table struct {
	Name    string
	Columns []*Column
	Indexes []*Index

	PrimaryKey string
}

func (t *Table) TableDefinition() string {
	s := fmt.Sprintf("CREATE TABLE %s (\n", t.Name)
	for i, c := range t.Columns {
		s += fmt.Sprintf("\t%s", c.ColumnDefinition())
		if i < len(t.Columns)-1 {
			s += ","

		}
		s += "\n"
	}
	if t.PrimaryKey != "" {
		s += fmt.Sprintf("\tPRIMARY KEY (%s)\n", t.PrimaryKey)
	}
	s += ");"

	return s
}

var tableNames = []string{
	"accounts",
	"customers",
	"items",
	"values",
	"orders", "stores", "order_items",
	"widgets", "categories", "keys", "objects", "tags",
	"roles", "permissions", "posts", "post_metadata",
	"metadata", "entries", "logs", "log_entries",
	"capabilities", "jobs", "queues", "job_logs",
	"queue_items",
}

// TODO(manuel) Create indexes
// TODO(manuel) Add all possible table options to test binlog and other potential weirdness

func GenerateTable(name string) *Table {
	hasId := rand.Intn(10) < 9
	primaryKey := ""

	usedColumns := make(map[string]bool)
	for _, column := range protoColumns {
		usedColumns[column.GetName()] = false
	}

	nColumns := randInRange(1, len(protoColumns))
	columns := []*Column{}
	if hasId {
		primaryKey = "id"
		columns = append(columns, &Column{
			Name:         "id",
			Nullable:     false,
			DatabaseType: "int",
			DefaultValue: "auto_increment",
			ProtoColumn:  &numericColumn{"id"},
		})
		usedColumns["id"] = true
	}

	for i := 0; i < nColumns; i++ {
		column := protoColumns[rand.Intn(len(protoColumns))]
		if usedColumns[column.GetName()] {
			i--
			continue
		}
		instance := column.Instance()
		if instance != nil {
			columns = append(columns, column.Instance())
			usedColumns[column.GetName()] = true
		}
	}

	if name == "" {
		name = randomString(tableNames)
	}

	return &Table{
		Name:       name,
		Columns:    columns,
		PrimaryKey: primaryKey,
	}
}
