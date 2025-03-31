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

var Application = newApplicationTable("public", "application", "")

type applicationTable struct {
	postgres.Table

	// Columns
	ID                   postgres.ColumnInteger
	Name                 postgres.ColumnString
	IapplicationAddress  postgres.ColumnString
	IconsensusAddress    postgres.ColumnString
	IinputboxAddress     postgres.ColumnString
	IinputboxBlock       postgres.ColumnFloat
	TemplateHash         postgres.ColumnString
	TemplateURI          postgres.ColumnString
	EpochLength          postgres.ColumnFloat
	DataAvailability     postgres.ColumnString
	State                postgres.ColumnString
	Reason               postgres.ColumnString
	LastInputCheckBlock  postgres.ColumnFloat
	LastOutputCheckBlock postgres.ColumnFloat
	ProcessedInputs      postgres.ColumnFloat
	CreatedAt            postgres.ColumnTimestampz
	UpdatedAt            postgres.ColumnTimestampz

	AllColumns     postgres.ColumnList
	MutableColumns postgres.ColumnList
}

type ApplicationTable struct {
	applicationTable

	EXCLUDED applicationTable
}

// AS creates new ApplicationTable with assigned alias
func (a ApplicationTable) AS(alias string) *ApplicationTable {
	return newApplicationTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new ApplicationTable with assigned schema name
func (a ApplicationTable) FromSchema(schemaName string) *ApplicationTable {
	return newApplicationTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new ApplicationTable with assigned table prefix
func (a ApplicationTable) WithPrefix(prefix string) *ApplicationTable {
	return newApplicationTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new ApplicationTable with assigned table suffix
func (a ApplicationTable) WithSuffix(suffix string) *ApplicationTable {
	return newApplicationTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newApplicationTable(schemaName, tableName, alias string) *ApplicationTable {
	return &ApplicationTable{
		applicationTable: newApplicationTableImpl(schemaName, tableName, alias),
		EXCLUDED:         newApplicationTableImpl("", "excluded", ""),
	}
}

func newApplicationTableImpl(schemaName, tableName, alias string) applicationTable {
	var (
		IDColumn                   = postgres.IntegerColumn("id")
		NameColumn                 = postgres.StringColumn("name")
		IapplicationAddressColumn  = postgres.StringColumn("iapplication_address")
		IconsensusAddressColumn    = postgres.StringColumn("iconsensus_address")
		IinputboxAddressColumn     = postgres.StringColumn("iinputbox_address")
		IinputboxBlockColumn       = postgres.FloatColumn("iinputbox_block")
		TemplateHashColumn         = postgres.StringColumn("template_hash")
		TemplateURIColumn          = postgres.StringColumn("template_uri")
		EpochLengthColumn          = postgres.FloatColumn("epoch_length")
		DataAvailabilityColumn     = postgres.StringColumn("data_availability")
		StateColumn                = postgres.StringColumn("state")
		ReasonColumn               = postgres.StringColumn("reason")
		LastInputCheckBlockColumn  = postgres.FloatColumn("last_input_check_block")
		LastOutputCheckBlockColumn = postgres.FloatColumn("last_output_check_block")
		ProcessedInputsColumn      = postgres.FloatColumn("processed_inputs")
		CreatedAtColumn            = postgres.TimestampzColumn("created_at")
		UpdatedAtColumn            = postgres.TimestampzColumn("updated_at")
		allColumns                 = postgres.ColumnList{IDColumn, NameColumn, IapplicationAddressColumn, IconsensusAddressColumn, IinputboxAddressColumn, IinputboxBlockColumn, TemplateHashColumn, TemplateURIColumn, EpochLengthColumn, DataAvailabilityColumn, StateColumn, ReasonColumn, LastInputCheckBlockColumn, LastOutputCheckBlockColumn, ProcessedInputsColumn, CreatedAtColumn, UpdatedAtColumn}
		mutableColumns             = postgres.ColumnList{NameColumn, IapplicationAddressColumn, IconsensusAddressColumn, IinputboxAddressColumn, IinputboxBlockColumn, TemplateHashColumn, TemplateURIColumn, EpochLengthColumn, DataAvailabilityColumn, StateColumn, ReasonColumn, LastInputCheckBlockColumn, LastOutputCheckBlockColumn, ProcessedInputsColumn, CreatedAtColumn, UpdatedAtColumn}
	)

	return applicationTable{
		Table: postgres.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:                   IDColumn,
		Name:                 NameColumn,
		IapplicationAddress:  IapplicationAddressColumn,
		IconsensusAddress:    IconsensusAddressColumn,
		IinputboxAddress:     IinputboxAddressColumn,
		IinputboxBlock:       IinputboxBlockColumn,
		TemplateHash:         TemplateHashColumn,
		TemplateURI:          TemplateURIColumn,
		EpochLength:          EpochLengthColumn,
		DataAvailability:     DataAvailabilityColumn,
		State:                StateColumn,
		Reason:               ReasonColumn,
		LastInputCheckBlock:  LastInputCheckBlockColumn,
		LastOutputCheckBlock: LastOutputCheckBlockColumn,
		ProcessedInputs:      ProcessedInputsColumn,
		CreatedAt:            CreatedAtColumn,
		UpdatedAt:            UpdatedAtColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
