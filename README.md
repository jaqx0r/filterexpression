# filterexpression

A parser for the [AIP-160 filter expression language](https://google.aip.dev/160) implemented in Go.

```go
import "github.com/jaqx0r/filterexpression"


   ...
   ast, err := filterexpression.Parse(req.filter)
   ...
```

TODO: an ast visitor to assist with query building.


