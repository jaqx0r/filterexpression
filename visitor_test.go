package filterexpression_test

import (
	"testing"

	"github.com/jaqx0r/filterexpression"
)

type TestVisitor struct {
	filterexpression.Visitor
}

func TestVisitFilter(t *testing.T) {
	ast, err := filterexpression.Parse(`book.title = "*and the*"`)
	if err != nil {
		t.Fatalf("Parse(): %v", err)
	}

	visitor := &TestVisitor{}
	if err := filterexpression.Visit(ast, visitor); err != nil {
		t.Errorf("Visit() failed: %v", err)
	}
}
