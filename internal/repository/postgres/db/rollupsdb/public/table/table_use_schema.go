//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

// UseSchema sets a new schema name for all generated table SQL builder types. It is recommended to invoke
// this method only once at the beginning of the program.
func UseSchema(schema string) {
	Application = Application.FromSchema(schema)
	Epoch = Epoch.FromSchema(schema)
	EspressoSchemaMigrations = EspressoSchemaMigrations.FromSchema(schema)
	ExecutionParameters = ExecutionParameters.FromSchema(schema)
	Input = Input.FromSchema(schema)
	NodeConfig = NodeConfig.FromSchema(schema)
	Output = Output.FromSchema(schema)
	Report = Report.FromSchema(schema)
	SchemaMigrations = SchemaMigrations.FromSchema(schema)
}
