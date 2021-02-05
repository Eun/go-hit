# yaegi-template [![Actions Status](https://github.com/Eun/yaegi-template/workflows/CI/badge.svg)](https://github.com/Eun/yaegi-template/actions) [![Coverage Status](https://coveralls.io/repos/github/Eun/yaegi-template/badge.svg)](https://coveralls.io/github/Eun/yaegi-template) [![PkgGoDev](https://img.shields.io/badge/pkg.go.dev-reference-blue)](https://pkg.go.dev/github.com/Eun/yaegi-template) [![GoDoc](https://godoc.org/github.com/Eun/yaegi-template?status.svg)](https://godoc.org/github.com/Eun/yaegi-template) [![go-report](https://goreportcard.com/badge/github.com/Eun/yaegi-template)](https://goreportcard.com/report/github.com/Eun/yaegi-template)
Use [yaegi](https://github.com/traefik/yaegi) as a template engine.

```go
package main

import (
	"os"

	"github.com/Eun/yaegi-template"
)

func main() {
	template := yaegi_template.MustNew(yaegi_template.DefaultOptions(), yaegi_template.DefaultSymbols()...)
	template.MustParseString(`
<html>
<$
    import (
        "fmt"
        "time"
    )
    func GreetUser(name string) {
        fmt.Printf("Hello %s, it is %s", name, time.Now().Format(time.Kitchen))
    }
$>

<p>
<$
    if context.LoggedIn {
        GreetUser(context.UserName)
    }
$>
</p>
</html>
`)

    type Context struct {
        LoggedIn bool
        UserName string
    }

    template.MustExec(os.Stdout, &Context{
        LoggedIn: true,
        UserName: "Joe Doe",
    })
}
```

## Example #2
You can use `<$-` to strip white spaces before the code block and
`-$>` to strip white spaces after the code block.  
Also omitting the print statement for simple evaluations is possible.
```go
package main

import (
	"os"

	"github.com/Eun/yaegi-template"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func main() {
	template := yaegi_template.MustNew(interp.Options{}, stdlib.Symbols)
	template.MustParseString(`
<html>
<p>
<$-
    context.UserName
-$>
</html>
`)

	type Context struct {
		LoggedIn bool
		UserName string
	}

	template.MustExec(os.Stdout, &Context{
		LoggedIn: true,
		UserName: "Joe Doe",
	})
}
```