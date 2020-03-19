# go-convert [![Travis](https://img.shields.io/travis/Eun/go-convert.svg)](https://travis-ci.org/Eun/go-convert) [![Codecov](https://img.shields.io/codecov/c/github/Eun/go-convert.svg)](https://codecov.io/gh/Eun/go-convert) [![GoDoc](https://godoc.org/github.com/Eun/go-convert?status.svg)](https://godoc.org/github.com/Eun/go-convert) [![go-report](https://goreportcard.com/badge/github.com/Eun/go-convert)](https://goreportcard.com/report/github.com/Eun/go-convert)
Convert a value into another type.

```bash
go get -u github.com/Eun/go-convert
```
## Usage
```go
// convert a int to a string
var s string
convert.MustConvert(1, &s)
fmt.Printf("%s\n", s)

// convert a map into a struct
type User struct {
	ID   int
	Name string
}
var u User
convert.MustConvert(map[string]string{
	"Name": "Joe",
	"ID":   "10",
}, &u)
fmt.Printf("%#v\n", u)

// convert Id to int and Groups to []int and keep the rest
m := map[string]interface{}{
	"Id":     0,
	"Groups": []int{},
}
// convert a map into well defined map
convert.MustConvert(
	map[string]interface{}{
		"Id":      "1",
		"Name":    "Joe",
		"Groups":  []string{"3", "6"},
		"Country": "US",
	},
	&m,
)
fmt.Printf("%v\n", m)

// convert a interface slice into well defined interface slice
// making the first one an integer, the second a string and the third an float
sl := []interface{}{0, "", 0.0}
convert.MustConvert([]string{"1", "2", "3"}, &sl)
fmt.Printf("%v\n", sl)
```

## Recipe system
_go-convert_ uses a recipe system that defines how and which types should be converted in which type.
A lot of recipes are already builtin (see [recipes.go](recipes.go)), however you can add your own or overwrite the builtin ones.
```go
type Roles struct {
    IsAdmin     bool
    IsDeveloper bool
}

type User struct {
    ID    int
    Name  string
    Roles Roles
}

// this is the data we want to convert
data := map[string]string{
    "id":    "10",
    "Name":  "Joe",
    "roles": "AD", // this user is Admin (A) and Developer (D)
}

// create a converter
conv := convert.New(convert.Options{
    Recipes: convert.MustMakeRecipes(
        // convert string into Roles
        func(_ convert.Converter, in string, out *Roles) error {
            (*out).IsAdmin = false
            (*out).IsDeveloper = false
            if strings.Contains(in, "A") {
                (*out).IsAdmin = true
            }
            if strings.Contains(in, "D") {
                (*out).IsDeveloper = true
            }
            return nil
        },
    ),
})

var user User
conv.MustConvert(data, &user)
// user is now an instance of User
fmt.Printf("%#v\n", user)

```




### Notice
This library is using reflection so be aware it might be slow in your usecase.  
