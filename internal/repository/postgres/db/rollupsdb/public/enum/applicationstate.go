//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var ApplicationState = &struct {
	Enabled    postgres.StringExpression
	Disabled   postgres.StringExpression
	Inoperable postgres.StringExpression
}{
	Enabled:    postgres.NewEnumValue("ENABLED"),
	Disabled:   postgres.NewEnumValue("DISABLED"),
	Inoperable: postgres.NewEnumValue("INOPERABLE"),
}
