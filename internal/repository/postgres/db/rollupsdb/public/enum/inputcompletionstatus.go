//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package enum

import "github.com/go-jet/jet/v2/postgres"

var InputCompletionStatus = &struct {
	None                       postgres.StringExpression
	Accepted                   postgres.StringExpression
	Rejected                   postgres.StringExpression
	Exception                  postgres.StringExpression
	MachineHalted              postgres.StringExpression
	OutputsLimitExceeded       postgres.StringExpression
	CycleLimitExceeded         postgres.StringExpression
	TimeLimitExceeded          postgres.StringExpression
	PayloadLengthLimitExceeded postgres.StringExpression
}{
	None:                       postgres.NewEnumValue("NONE"),
	Accepted:                   postgres.NewEnumValue("ACCEPTED"),
	Rejected:                   postgres.NewEnumValue("REJECTED"),
	Exception:                  postgres.NewEnumValue("EXCEPTION"),
	MachineHalted:              postgres.NewEnumValue("MACHINE_HALTED"),
	OutputsLimitExceeded:       postgres.NewEnumValue("OUTPUTS_LIMIT_EXCEEDED"),
	CycleLimitExceeded:         postgres.NewEnumValue("CYCLE_LIMIT_EXCEEDED"),
	TimeLimitExceeded:          postgres.NewEnumValue("TIME_LIMIT_EXCEEDED"),
	PayloadLengthLimitExceeded: postgres.NewEnumValue("PAYLOAD_LENGTH_LIMIT_EXCEEDED"),
}
