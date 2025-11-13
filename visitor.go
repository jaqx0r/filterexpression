package filterexpression

// FilterVisitor provides an interface for visiting a Filter expression AST.
// Interface functions are called as each are visited.  Each method can return
// an error to indicate a construction failure or a semantic error, which
// immediately halts the visitor, returning that error to the caller of the
// Visit function below.  If the function returns nil the Visit method will continue to the next node in the ST.
type FilterVisitor interface {
	// VisitSequence inspects a sequence of factors, each equal and mandatory
	// in the search.  A sequence is equivalent to a logical AND with exact
	// match semantics.
	VisitSequence(ast *Sequence) error

	// VisitFactor inspects a disjunction of terms, logically ORed together.
	VisitFactor(ast *Factor) error

	// VisitTerm visits a unary expression, possibly negated.
	VisitTerm(ast *Term) error

	// VisitRestriction inpects a Restriction production, which describes a
	// comparison relation.
	VisitRestriction(ast *Restriction) error

	// VisitFunction is called when visiting a Function production which describes a function call and arguments.
	VisitFunction(ast *Function) error

	// VisitMember is called when visiting a Member production which describes a dot-qualfied field reference.
	VisitMember(ast *Member) error
}

// A base Visitor that satisfies the FilterVisitor interface.
// Embed this struct into your own Visitor so you only need to implement the methods you require.
type Visitor struct {
}

func (Visitor) VisitSequence(ast *Sequence) error {
	return nil
}

func (Visitor) VisitFactor(ast *Factor) error {
	return nil
}

func (Visitor) VisitTerm(ast *Term) error {
	return nil
}

func (Visitor) VisitRestriction(ast *Restriction) error {
	return nil
}

func (Visitor) VisitFunction(ast *Function) error {
	return nil
}

func (Visitor) VisitMember(ast *Member) error {
	return nil
}

func Visit(ast *Filter, visitor FilterVisitor) error {
	return ast.Accept(visitor)
}

func (ast *Filter) Accept(visitor FilterVisitor) error {
	for _, e := range ast.Expression {
		err := e.Accept(visitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ast *Expression) Accept(visitor FilterVisitor) error {
	for _, s := range ast.Sequence {
		err := s.Accept(visitor)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ast *Sequence) Accept(visitor FilterVisitor) error {
	for _, f := range ast.Factor {
		err := f.Accept(visitor)
		if err != nil {
			return err
		}
	}
	return visitor.VisitSequence(ast)
}

func (ast *Factor) Accept(visitor FilterVisitor) error {
	for _, t := range ast.Term {
		err := t.Accept(visitor)
		if err != nil {
			return err
		}
	}
	return visitor.VisitFactor(ast)
}

func (ast *Term) Accept(visitor FilterVisitor) error {
	err:= ast.Simple.Accept(visitor)
	if err != nil {
		return err
	}
	return visitor.VisitTerm(ast)
}

func (ast *Simple) Accept(visitor FilterVisitor) error {
	if ast.Restriction != nil {
		return ast.Restriction.Accept(visitor)
	}
	if ast.Composite != nil {
		return ast.Composite.Accept(visitor)
	}
	return nil
}

func (ast *Restriction) Accept(visitor FilterVisitor) error {
	err := ast.Comparable.Accept(visitor)
	if err != nil {
		return err
	}
	err = ast.Arg.Accept(visitor)
	if err != nil {
		return err
	}
	return visitor.VisitRestriction(ast)
}

func (ast *Comparable) Accept(visitor FilterVisitor) error {
	if ast.Function != nil {
		err := ast.Function.Accept(visitor)
		if err != nil {
			return err
		}
		return visitor.VisitFunction(ast.Function)
	}
	if ast.Member != nil {
		err := ast.Member.Accept(visitor)
		if err != nil {
			return err
		}
		return visitor.VisitMember(ast.Member)
	}
	return nil
}

func (ast *Arg) Accept(visitor FilterVisitor) error {
	return nil
}

func (ast *Function) Accept(visitor FilterVisitor) error {
	return nil
}

func (ast *Member)Accept(visitor FilterVisitor) error {
	return nil
}
