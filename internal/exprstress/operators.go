package exprstress

import (
	"math/rand"

	"github.com/zclconf/go-cty/cty"
)

type arithmeticExpr struct {
	LHS      Expression
	RHS      Expression
	Operator string
}

var arithmeticOperators = []string{
	"+", "-", "*", "/", "%",
}

func generateArithmeticOperator(depth int) expressionGenerator {
	generateOperand := generateArithmeticOperand(depth + 1)
	return func(rand *rand.Rand) Expression {
		n := rand.Intn(len(arithmeticOperators))
		op := arithmeticOperators[n]
		lhs := generateOperand(rand)
		rhs := generateOperand(rand)
		return &arithmeticExpr{
			LHS:      lhs,
			RHS:      rhs,
			Operator: op,
		}
	}
}

func (e arithmeticExpr) BuildSource(w SourceWriter) {
	// To ensure we get the expected evaluation order without having to
	// analyze for precedence, we'll enclose both operands in parentheses.
	// This does mean that we're not actually testing precedence rules here,
	// but that's okay because HCL has its own tests for that.
	w.WriteString("(")
	e.LHS.BuildSource(w)
	w.WriteString(") ")
	w.WriteString(e.Operator)
	w.WriteString(" (")
	e.RHS.BuildSource(w)
	w.WriteString(")")
}

func (e arithmeticExpr) ExpectedResult() Expected {
	mode := SpecifiedValue
	lhsExpect := e.LHS.ExpectedResult()
	rhsExpect := e.RHS.ExpectedResult()
	if lhsExpect.Mode == UnknownValue || rhsExpect.Mode == UnknownValue {
		mode = UnknownValue
	}
	sensitive := lhsExpect.Sensitive || rhsExpect.Sensitive

	return Expected{
		Type:      cty.Number,
		Mode:      mode,
		Sensitive: sensitive,
	}
}
