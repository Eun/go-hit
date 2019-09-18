# go-convert

Convert a value into another type.

```bash
go get -u github.com/Eun/go-convert
```
## Usage
```go
// convert a int to a string
fmt.Printf("%v\n", convert.MustConvert(1, ""))


// convert a map into a struct
type User struct {
	ID int
	Name string
}
fmt.Printf("%#v\n", convert.MustConvert(
	map[string]string{
	    "Name": "Joe",
	    "ID": "10",
	},
	User{}))

// convert a map into well defined map
fmt.Printf("%v\n", MustConvert(
	map[string]interface{}{
		"Id":        "1",
		"Name":      "Joe",
		"Groups":    []string{"3", "6"},
		"Country":   "US",
 	},
    // convert Id to int and Groups to []int and keep the rest
	map[string]interface{}{
        "Id":      0,
		"Groups":  []int{},
 	}, 
))

// convert a interface slice into well defined interface slice
// making the first one an integer, the second a string and the third an float
fmt.Printf("%v\n", MustConvert([]string{"1", "2", "3"}, []interface{}{0, "", 0.0}))
```


### Notice
This library is using reflection so be aware it might be slow in your usecase.  
