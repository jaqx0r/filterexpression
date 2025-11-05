package filterexpression_test

import (
	"testing"

	participle "github.com/alecthomas/participle/v2"
	"github.com/jaqx0r/filterexpression"
	"github.com/kr/pretty"
)

type Test[F any] struct {
	name  string
	input string
}

func (tf Test[T]) Test(t *testing.T) {
	typeName := *new(T)
	parser, err := participle.Build[T](
		participle.Lexer(filterexpression.Lexer),
		participle.Elide("Whitespace"),
		participle.UseLookahead(4),
	)
	if err != nil {
		t.Errorf("participle.Build[%v]() failed: %v", typeName, err)
	}
	t.Logf("parser ebnf:\n%s", parser.String())

	ast, err := parser.ParseString(tf.name, tf.input)
	if err != nil {
		t.Errorf("parser.ParseString(%q) failed: %v", tf.input, err)
		t.Logf("parse tree so far: \n%+# v", pretty.Formatter(ast))
	}
}

func (tf Test[T]) Name() string {
	return tf.name
}

type Testable interface {
	Test(t *testing.T)
	Name() string
}

func TestParsers(t *testing.T) {
	for _, tc := range []Testable{
		Test[filterexpression.Name]{"name", "name"},
		Test[filterexpression.Name]{"name is keyword", "AND"},

		Test[filterexpression.Field]{"field", "field"},
		Test[filterexpression.Field]{"field is keyword", "OR"},

		Test[filterexpression.Value]{"value", "value"},
		Test[filterexpression.Value]{"value is string", "\"value\""},
		Test[filterexpression.Value]{"value is string with asterisks", "\"*value*\""},

		Test[filterexpression.Function]{"function no args", "func()"},
		Test[filterexpression.Function]{"function one arg", "func(a)"},
		Test[filterexpression.Function]{"nested function one arg", "func.func(a)"},
		Test[filterexpression.Function]{"nested function two args", "func.func(a, b)"},
		Test[filterexpression.Function]{"function example 1", "regex(m.key, '^.*prod.*$')"},
		Test[filterexpression.Function]{"function example 2", "math.mem('30mb')"},

		Test[filterexpression.Member]{"Member", "expr.type_map.1.type"},

		Test[filterexpression.Composite]{"composite example", `(msg.endsWith('world') AND retries < 10)`},

		Test[filterexpression.Restriction]{"restriction example equality", `package=com.google`},
		Test[filterexpression.Restriction]{"restriction example inequality", `msg != hello`},
		Test[filterexpression.Restriction]{"restriction example greater than", `1 > 0`},
		Test[filterexpression.Restriction]{"restriction example greater or equal", `2.5 >= 2.4`},
		Test[filterexpression.Restriction]{"restriction example less than", `yesterday < request.time`},
		Test[filterexpression.Restriction]{"restriction example less or equal", `experiment.rollout <= cohort(request.user)`},
		Test[filterexpression.Restriction]{"restriction example has", `map:key`},
		Test[filterexpression.Restriction]{"restriction example global", `prod`},

		Test[filterexpression.Term]{"term example logical not", "NOT (a OR b)"},
		Test[filterexpression.Term]{"term example alternative not", `-file:".java"`},
		Test[filterexpression.Term]{"term example negation", `-30`},

		Test[filterexpression.Factor]{"factor example", "a < 10 OR a >= 100"},
		Test[filterexpression.Sequence]{"sequence example", "New York Giants OR Yankees"},

		Test[filterexpression.Expression]{"expression example", "a b AND c AND d"},
	} {
		t.Run(tc.Name(),
			tc.Test)
	}
}
