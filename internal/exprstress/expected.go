package exprstress

import (
	"github.com/zclconf/go-cty/cty"
)

// Expected represents some cross-cutting metadata about an expected expression
// result, which we use both to allow intermediate expressions to make
// expectations about their own results based expectations of their inputs
// and also to verify that the result of an overall test expression matches
// the final expectations.
type Expected struct {
	// Type is a cty type that the final result type must match exactly.
	// (This is not a type _constraint_, so dynamic pseudo-type may appear
	// within it only if the expected result will be unknown, null, or an empty
	// collection.)
	Type cty.Type

	// Mode indicates whether the result is expected to be unknown, null,
	// or neither.
	Mode ValueMode

	// Sensitive indicates whether the result is expected to be marked as
	// sensitive.
	Sensitive bool
}

// CouldConvertTo returns true if the reciever describes a value that
// definitely could become a known value of the given type under type
// conversion.
//
// A return value of false doesn't mean that such a conversion would
// fail, but only that we can't statically prove that it would succeed.
func (e Expected) CouldConvertTo(ty cty.Type) bool {
	switch e.Mode {
	case UnknownValue, NullValue:
		if e.Type == cty.DynamicPseudoType {
			// A null or unknown value of DynamicPseudoType can convert
			// to a null or unknown value of any other type.
			return true
		}
	}

	switch {
	case ty == cty.String:
		switch e.Type {
		case cty.String, cty.Number, cty.Bool:
			return true
		default:
			return false
		}
	default:
		return e.Type.Equals(ty)
	}
}

// ValueMode represents the three mutually-exclusive modes a value can be in:
// unknown, null, or known-and-not-null ("specified").
type ValueMode rune

//go:generate go run golang.org/x/tools/cmd/stringer -type=ValueMode -output=value_mode_string.go expected.go

const (
	// SpecifiedValue represents a value that is known and not null.
	SpecifiedValue ValueMode = 'C'

	// UnknownValue represents an unknown value.
	UnknownValue ValueMode = 'U'

	// NullValue represents a known null value.
	NullValue ValueMode = 'N'
)

// GoString implements fmt.GoStringer.
func (m ValueMode) GoString() string {
	return "exprstress." + m.String()
}
