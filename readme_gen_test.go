// +build doctest

package hit_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	// ⚠️⚠️⚠️ This file was autogenerated by generators/readme/readme ⚠️⚠️⚠️ //
	"net/http"

	. "github.com/Eun/go-hit"
)

func TestReadmeCodePart0(t *testing.T) {
	MustDo(
		Get("https://httpbin.org/post"),
		Expect().Status().Equal(http.StatusMethodNotAllowed),
		Expect().Headers("Content-Type").NotEmpty(),
		Expect().Body().String().Contains("Method Not Allowed"),
	)
}
func TestReadmeCodePart1(t *testing.T) {
	MustDo(
		Post("https://httpbin.org/post"),
		Send().Body().String("Hello HttpBin"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains("Hello HttpBin"),
	)
}
func TestReadmeCodePart2(t *testing.T) {
	MustDo(
		Post("https://httpbin.org/post"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string][]string{"Foo": []string{"Bar", "Baz"}}),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().JSON().JQ(".json.Foo[1]").Equal("Baz"),
	)
}
func TestReadmeCodePart3(t *testing.T) {
	var name string
	var roles []string
	MustDo(
		Post("https://httpbin.org/post"),
		Send().Headers("Content-Type").Add("application/json"),
		Send().Body().JSON(map[string]interface{}{"Name": "Joe", "Roles": []string{"Admin", "Developer"}}),
		Expect().Status().Equal(http.StatusOK),
		Store().Response().Body().JSON().JQ(".json.Name").In(&name),
		Store().Response().Body().JSON().JQ(".json.Roles").In(&roles),
	)
	fmt.Printf("%s has %d roles\n", name, len(roles))
}
func TestReadmeCodePart4(t *testing.T) {
	MustDo(
		Post("https://httpbin.org/post"),
		Debug(),
		Debug().Response().Body(),
	)
}
func TestReadmeCodePart5(t *testing.T) {
	MustDo(
		Post("https://httpbin.org/post"),
		Expect().Status().Equal(200),
		Send().Body().String("Hello World"),
		Expect().Body().String().Contains("Hello"),
	)
}
func TestReadmeCodePart6(t *testing.T) {
	MustDo(
		Get("https://httpbin.org/get"),
		Send().Custom(func(hit Hit) {
			hit.Request().Body().SetStringf("Hello %s", "World")
		}),
		Expect().Custom(func(hit Hit) {
			if len(hit.Response().Body().MustString()) <= 0 {
				t.FailNow()
			}
		}),
		Custom(AfterExpectStep, func(Hit) {
			fmt.Println("everything done")
		}),
	)
}
func TestReadmeCodePart7(t *testing.T) {
	template := CombineSteps(
		Post("https://httpbin.org/post"),
		Send().Headers("Content-Type").Add("application/json"),
		Expect().Headers("Content-Type").Equal("application/json"),
	)
	MustDo(
		template,
		Send().Body().JSON("Hello World"),
	)

	MustDo(
		template,
		Send().Body().JSON("Hello Universe"),
	)
}
func TestReadmeCodePart8(t *testing.T) {
	template := CombineSteps(
		Get("https://httpbin.org/basic-auth/joe/secret"),
		Expect().Status().Equal(http.StatusOK),
	)
	MustDo(
		Description("login with correct credentials"),
		template,
		Send().Headers("Authorization").Add("Basic "+base64.StdEncoding.EncodeToString([]byte("joe:secret"))),
	)

	Test(t,
		Description("login with incorrect credentials"),
		template,
		Clear().Expect().Status(),
		Expect().Status().Equal(http.StatusUnauthorized),
		Send().Headers("Authorization").Add("Basic "+base64.StdEncoding.EncodeToString([]byte("joe:joe"))),
	)
}
