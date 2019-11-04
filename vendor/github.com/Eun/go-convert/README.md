# go-convert

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


### Notice
This library is using reflection so be aware it might be slow in your usecase.  
