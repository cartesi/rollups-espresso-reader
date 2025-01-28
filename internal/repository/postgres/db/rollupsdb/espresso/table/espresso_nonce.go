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

var EspressoNonce = newEspressoNonceTable("espresso", "espresso_nonce", "")

type espressoNonceTable struct {
	postgres.Table

	// Columns
	SenderAddress      postgres.ColumnString
	ApplicationAddress postgres.ColumnString
	Nonce              postgres.ColumnInteger

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type EspressoNonceTable struct {
	espressoNonceTable

	EXCLUDED espressoNonceTable
}

// AS creates new EspressoNonceTable with assigned alias
func (a EspressoNonceTable) AS(alias string) *EspressoNonceTable {
	return newEspressoNonceTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new EspressoNonceTable with assigned schema name
func (a EspressoNonceTable) FromSchema(schemaName string) *EspressoNonceTable {
	return newEspressoNonceTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new EspressoNonceTable with assigned table prefix
func (a EspressoNonceTable) WithPrefix(prefix string) *EspressoNonceTable {
	return newEspressoNonceTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new EspressoNonceTable with assigned table suffix
func (a EspressoNonceTable) WithSuffix(suffix string) *EspressoNonceTable {
	return newEspressoNonceTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newEspressoNonceTable(schemaName, tableName, alias string) *EspressoNonceTable {
	return &EspressoNonceTable{
		espressoNonceTable: newEspressoNonceTableImpl(schemaName, tableName, alias),
		EXCLUDED:           newEspressoNonceTableImpl("", "excluded", ""),
	}
}

func newEspressoNonceTableImpl(schemaName, tableName, alias string) espressoNonceTable {
	var (
		SenderAddressColumn      = postgres.StringColumn("sender_address")
		ApplicationAddressColumn = postgres.StringColumn("application_address")
		NonceColumn              = postgres.IntegerColumn("nonce")
		allColumns               = postgres.ColumnList{SenderAddressColumn, ApplicationAddressColumn, NonceColumn}
		mutableColumns           = postgres.ColumnList{SenderAddressColumn, ApplicationAddressColumn, NonceColumn}
	)

	return espressoNonceTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		SenderAddress:      SenderAddressColumn,
		ApplicationAddress: ApplicationAddressColumn,
		Nonce:              NonceColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
