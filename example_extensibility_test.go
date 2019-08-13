package hit_test

// import (
// 	"fmt"
//
// 	"strings"
//
// 	"github.com/Eun/go-hit"
// )
//
// type User struct {
// 	ID   int
// 	Name string
// }
//
// // MyTestingFramework embeds Hit so we can add new/replace functions to/from Hit
// type MyTestingFramework struct {
// 	hit.Hit
// }
//
// // Add a complete new function
// func (my MyTestingFramework) CheckTheSparkles() MyTestingFramework {
// 	fmt.Println("Lets see if we can find sparkles ✨✨✨...")
// 	my.Expect().Custom(func(hit hit.Hit) {
// 		// use your custom code
// 		if !strings.Contains(hit.Response().Body().String(), "✨") {
// 			hit.T().Errorf("No sparkles present, too sad :-(")
// 		}
// 		// or use our expect
// 		// hit.IExpect().Body().Contains("✨")
// 	})
// 	return my
// }
//
// // Overwrite Hit's Send function and return with our custom Send struct,
// // now we can implement a custom Send function
// func (my MyTestingFramework) Send(data ...interface{}) MySend {
// 	return MySend{
// 		Send: my.Hit.Send(data...),
// 	}
// }
//
// type MySend struct {
// 	hit.Send
// }
//
// // User sets the body to the user
// func (snd MySend) User(user User) MyTestingFramework {
// 	// We have to return our Framework to make sure we can use our custom functions later on
// 	return MyTestingFramework{
// 		Hit: snd.Custom(func(hit hit.Hit) {
// 			// use hit.Send()
// 			hit.Send().Headers().Set("Content-Type", "application/json").
// 				Send().Body().JSON(user)
//
// 			// or hit.Request()
// 			// hit.Request().Header.Set("Content-Type", "application/json")
// 			// hit.Request().Body().JSON().Set(user)
// 		}),
// 	}
// }
//
// // Overwrite Hit's IExpect function and return with our custom IExpect struct,
// // now we can implement a custom Send function
// func (my MyTestingFramework) Expect(data ...interface{}) MyExpect {
// 	return MyExpect{
// 		IExpect: my.Hit.Expect(data...),
// 	}
// }
//
// type MyExpect struct {
// 	hit.IExpect
// }
//
// // User expects the body to be equal to the specific user
// func (exp MyExpect) User(user User) MyTestingFramework {
// 	// We have to return our Framework to make sure we can use our custom functions later on
// 	return MyTestingFramework{
// 		Hit: exp.Custom(func(hit hit.Hit) {
// 			// use hit.IExpect
// 			hit.Expect().Headers("Content-Type").Equal("application/json").
// 				Expect().Body().JSON().Equal("", user)
//
// 			// or hit.Response()
// 			// if hit.Response().Header.Get("Content-Type") != "application/json" {
// 			// 	hit.T().Errorf("%#v != %#v", hit.Response().Header.Get("Content-Type"), "application/json")
// 			// }
// 			// if !reflect.DeepEqual(hit.Response().Body().JSON().GetAs(&User{}), &user) {
// 			// 	hit.T().Errorf("%#v != %#v", hit.Response().Body().JSON().GetAs(&User{}), &user)
// 			// }
// 		}),
// 	}
// }
//
// func Example_extensibility() {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	// Create a new instance of MyTestingFramework
// 	// and do all steps on this Instance
// 	(MyTestingFramework{
// 		Hit: hit.Post(hit.PanicT{}, s.URL),
// 	}).
// 		Send().User(User{10, "Joe ✨"}).
// 		Expect().User(User{10, "Joe ✨"}).
// 		CheckTheSparkles().
// 		Do()
// }
