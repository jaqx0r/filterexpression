// Package filterexpression parses List filter expressions defined in [AIP-160](https://google.aip.dev/160)
//
// The full specification in EBNF is https://google.aip.dev/assets/misc/ebnf-filtering.txt
package filterexpression

import (
	participle "github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Filter, possibly empty
type Filter struct {
	Pos lexer.Position

	Expression []Expression `@@*`
}

// Expression may either be a conjunction (AND) of sequences or a simple sequence.
//
// Note the AND is case sensitive.
//
// Example: `a b AND c AND d`
//
// The expression `(a b) AND c AND d` is equivalent to the example.
type Expression struct {
	Pos lexer.Position

	Sequence []Sequence `@@ ( "AND" @@ )*`
}

// Sequence is composed of one or more whitespace (WS) separated factors.
//
// A sequence expresses a logical relationship between 'factors' where
// the ranking of a filter result may be scored according to the number
// factors that match and other such criteria as the proximity of factors
// to each other within a document.
//
// When filters are used with exact match semantics rather than fuzzy
// match semantics, a sequence is equivalent to AND.
//
// Example: `New York Giants OR Yankees`
//
// The expression `New York (Giants OR Yankees)` is equivalent to the
// example.
type Sequence struct {
	Pos lexer.Position

	Factor []Factor `@@+`
}

// Factors may either be a disjunction (OR) of terms or a simple term.
//
// Note, the OR is case-sensitive.
//
// Example: `a < 10 OR a >= 100`
type Factor struct {
	Pos lexer.Position

	Term []Term `@@ ( "OR" @@ )*`
}

// Terms may either be unary or simple expressions.
//
// Unary expressions negate the simple expression, either mathematically `-`
// or logically `NOT`. The negation styles may be used interchangeably.
//
// Note, the `NOT` is case-sensitive and must be followed by at least one
// whitespace (WS).
//
// Examples:
// * logical not     : `NOT (a OR b)`
// * alternative not : `-file:".java"`
// * negation        : `-30`
type Term struct {
	Pos lexer.Position

	Negate bool `@( "NOT" | "-" )?`

	Simple Simple `@@`
}

// Simple expressions may either be a restriction or a nested (composite)
// expression.
type Simple struct {
	Pos lexer.Position

	Restriction Restriction `@@`
	Composite   Composite   `| @@`
}

// Restrictions express a relationship between a comparable value and a
// single argument. When the restriction only specifies a comparable
// without an operator, this is a global restriction.
//
// Note, restrictions are not whitespace sensitive.
//
// Examples:
// * equality         : `package=com.google`
// * inequality       : `msg != 'hello'`
// * greater than     : `1 > 0`
// * greater or equal : `2.5 >= 2.4`
// * less than        : `yesterday < request.time`
// * less or equal    : `experiment.rollout <= cohort(request.user)`
// * has              : `map:key`
// * global           : `prod`
//
// In addition to the global, equality, and ordering operators, filters
// also support the has (`:`) operator. The has operator is unique in
// that it can test for presence or value based on the proto3 type of
// the `comparable` value. The has operator is useful for validating the
// structure and contents of complex values.
type Restriction struct {
	Pos lexer.Position

	Comparable Comparable `@@`
	Comparator Comparator `( @( "<=" | "<" | ">=" | ">" | "!=" | "=" | ":" )`
	Arg        Arg        ` @@ )?`
}

// Comparable may either be a member or function.
type Comparable struct {
	Pos lexer.Position

	Function Function `@@`
	Member   Member   `| @@`
}

// Member expressions are either value or DOT qualified field references.
//
// Example: `expr.type_map.1.type`
type Member struct {
	Pos lexer.Position

	Value  Value   `@@`
	Fields []Field `( "." @@ )*`
}

// Function calls may use simple or qualified names with zero or more
// arguments.
//
// All functions declared within the list filter, apart from the special
// `arguments` function must be provided by the host service.
//
// Examples:
// * `regex(m.key, '^.*prod.*$')`
// * `math.mem('30mb')`
//
// Antipattern: simple and qualified function names may include keywords:
// NOT, AND, OR. It is not recommended that any of these names be used
// within functions exposed by a service that supports list filters.
type Function struct {
	Pos lexer.Position

	Name []Name `@@ ( "." @@ )*`
	Args []Arg  `"(" ( @@ ( "," @@ )* )? ")"`
}

type Comparator int

const (
	CompLessEquals Comparator = iota
	CompLessThan
	CompGreaterEquals
	CompGreaterThan
	CompNotEquals
	CompEquals
	CompHas
)

var comparatorMap = map[string]Comparator{
	"<=": CompLessEquals,
	"<":  CompLessThan,
	">=": CompGreaterEquals,
	">":  CompGreaterThan,
	"!=": CompNotEquals,
	"=":  CompEquals,
	":":  CompHas,
}

func (c *Comparator) Capture(s []string) error {
	*c = comparatorMap[s[0]]
	return nil
}

// Composite is a parenthesized expression, commonly used to group
// terms or clarify operator precedence.
//
// Example: `(msg.endsWith('world') AND retries < 10)`
type Composite struct {
	Pos lexer.Position

	Expression Expression `"(" @@ ")"`
}

// Value may either be a TEXT or STRING.
//
// TEXT is a free-form set of characters without whitespace (WS)
// or . (DOT) within it. The text may represent a variable, string,
// number, boolean, or alternative literal value and must be handled
// in a manner consistent with the service's intention.
//
// STRING is a quoted string which may or may not contain a special
// wildcard `*` character at the beginning or end of the string to
// indicate a prefix or suffix-based search within a restriction.
type Value struct {
	Pos lexer.Position

	Text   string `@Text`
	String string `| @String`
}

type Field struct {
	Pos lexer.Position

	Value   *Value `@@`
	Keyword string `| @( "AND" | "OR" | "NOT" )`
}

type Arg struct {
	Pos lexer.Position

	Comparable Comparable `@@`
	Composite  Composite  `| @@`
}

type Name struct {
	Pos lexer.Position

	Text    string `@Text`
	Keyword string `| @( "AND" | "OR" | "NOT" )`
}

var Lexer = lexer.MustSimple([]lexer.SimpleRule{
	{Name: "Text", Pattern: `[a-zA-Z0-9_]+`},
	{Name: "String", Pattern: `['"]\*?(\\'|\\"|[^"'])*\*?['"]`},
	{Name: "Keyword", Pattern: `\b(AND|OR|NOT)\b`},
	{Name: "Whitespace", Pattern: `[ \t\n\r]+`},
	{Name: "Operators", Pattern: `<=|>=|!=|[=\:.<>=(),-]`},
})

var DefaultParser = participle.MustBuild[Filter](
	participle.Lexer(Lexer),
	participle.Elide("Whitespace"),
	// 7 is an arbitrary number that lets us fall back to parse Member
	// instead of Function after 4 Value tokens (7 including dots.)
	// Unfortunately inverting the production will never match a Function
	// because the parser is non-greedy.
	participle.UseLookahead(7),
)

// Parse parses the given expression into a Filter AST.
//
// If the expression is not compliant with [AIP-160](https://google.aip.dev/160) a parse error is raised and a best effort parse tree is returned.
func Parse(expression string) (*Filter, error) {
	return DefaultParser.ParseString("", expression)
}
