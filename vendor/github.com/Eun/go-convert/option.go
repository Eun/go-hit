package convert

// Options can be used to alter the behavior of the converter
type Options struct {
	// SkipUnknownFields can be used to ignore fields that are not existent in the destination type
	//
	// Example:
	//     type User struct {
	//         Name string
	//     }
	//     m := map[string]interface{}{
	//         "Name":     "Joe",
	//         "Surname":  "Doe",
	//     }
	//     // convert a map into User type
	//     var user User
	//     // will fail because Surname is not in the User struct
	//     MustConvert(m, &user)
	//     // will pass
	//     MustConvert(m, &user, Options{SkipUnknownFields: true})
	SkipUnknownFields bool
	// Recipes can be used to define custom recipes for this converter
	Recipes []Recipe
}
