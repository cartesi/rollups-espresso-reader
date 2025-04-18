//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/postgres"
)

var AppInfo = newAppInfoTable("espresso", "app_info", "")

type appInfoTable struct {
	postgres.Table

	// Columns
	ApplicationAddress         postgres.ColumnString
	StartingBlock              postgres.ColumnFloat
	Namespace                  postgres.ColumnFloat
	LastProcessedEspressoBlock postgres.ColumnFloat
	Index                      postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type AppInfoTable struct {
	appInfoTable

	EXCLUDED appInfoTable
}

// AS creates new AppInfoTable with assigned alias
func (a AppInfoTable) AS(alias string) *AppInfoTable {
	return newAppInfoTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new AppInfoTable with assigned schema name
func (a AppInfoTable) FromSchema(schemaName string) *AppInfoTable {
	return newAppInfoTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new AppInfoTable with assigned table prefix
func (a AppInfoTable) WithPrefix(prefix string) *AppInfoTable {
	return newAppInfoTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new AppInfoTable with assigned table suffix
func (a AppInfoTable) WithSuffix(suffix string) *AppInfoTable {
	return newAppInfoTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newAppInfoTable(schemaName, tableName, alias string) *AppInfoTable {
	return &AppInfoTable{
		appInfoTable: newAppInfoTableImpl(schemaName, tableName, alias),
		EXCLUDED:     newAppInfoTableImpl("", "excluded", ""),
	}
}

func newAppInfoTableImpl(schemaName, tableName, alias string) appInfoTable {
	var (
		ApplicationAddressColumn         = postgres.StringColumn("application_address")
		StartingBlockColumn              = postgres.FloatColumn("starting_block")
		NamespaceColumn                  = postgres.FloatColumn("namespace")
		LastProcessedEspressoBlockColumn = postgres.FloatColumn("last_processed_espresso_block")
		IndexColumn                      = postgres.IntegerColumn("index")
		allColumns                       = postgres.ColumnList{ApplicationAddressColumn, StartingBlockColumn, NamespaceColumn, LastProcessedEspressoBlockColumn, IndexColumn}
		mutableColumns                   = postgres.ColumnList{StartingBlockColumn, NamespaceColumn, LastProcessedEspressoBlockColumn, IndexColumn}
	)

	return appInfoTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ApplicationAddress:         ApplicationAddressColumn,
		StartingBlock:              StartingBlockColumn,
		Namespace:                  NamespaceColumn,
		LastProcessedEspressoBlock: LastProcessedEspressoBlockColumn,
		Index:                      IndexColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
