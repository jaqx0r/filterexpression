# filterexpression

A parser for the [AIP-160 filter expression language](https://google.aip.dev/160) implemented in Go.

```go
import "github.com/jaqx0r/filterexpression"


   ...
   ast, err := filterexpression.Parse(req.filter)
   ...
```

Visit the AST by implementing the `FilterVisitor` interface.

You can embed the existing `Visitor` as a base, so your implementation only needs to override the methods it cares about.

```go
type Visitor struct {
    filtervisitor.Visitor

    query query.Builder
}

func (v *Visitor) VisitFunction(ast *filterexpression.Function) error {
  query.Function(ast.Name[0])
  return nil
}

...
   visitor = &Visitor{}
   if err := filterexpression.Visit(ast, visitor); err != nil {
     log.Errorf("Visit() failed: %v", err)
   }
```
