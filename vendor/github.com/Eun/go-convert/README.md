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
```

### Notice
This library is using reflection so be aware it might be slow in your usecase.  
