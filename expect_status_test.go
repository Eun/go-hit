package hit_test

//
// func TestExpectStatus_Equal(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status(http.StatusOK).
// 			Do()
// 	})
//
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("Expected status code to be 404 but was 200 instead")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status(http.StatusNotFound).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatusOneOf(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().OneOf(200, 300).
// 			Do()
// 	})
//
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr(`[]int{`), PtrStr(`300,`), PtrStr(`400,`), PtrStr(`} does not contain 200`)), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().OneOf(300, 400).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatus_GreaterThan(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().GreaterThan(199).
// 			Do()
// 	})
//
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("expected 200 to be greater than 299")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().GreaterThan(299).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatus_GreaterOrEqualThan(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().GreaterOrEqualThan(200).
// 			Do()
// 	})
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("expected 200 to be greater or equal than 299")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().GreaterOrEqualThan(299).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatus_LessThan(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().LessThan(201).
// 			Do()
// 	})
//
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("expected 200 to be less than 100")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().LessThan(100).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatus_LessOrEqualThan(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().LessOrEqualThan(200).
// 			Do()
// 	})
//
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("expected 200 to be less or equal than 100")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().LessOrEqualThan(100).
// 				Do()
// 		})
// 	})
// }
//
// func TestExpectStatus_Between(t *testing.T) {
// 	s := EchoServer()
// 	defer s.Close()
//
// 	t.Run("ok", func(t *testing.T) {
// 		hit.Post(t, s.URL).
// 			Send().Body("Hello World").
// 			Expect().Status().Between(100, 200).
// 			Do()
// 	})
// 	t.Run("failing", func(t *testing.T) {
// 		require.Panics(t, func() {
// 			hit.Post(NewPanicWithMessage(t, PtrStr("expected 200 to be between 300 and 400")), s.URL).
// 				Send().Body("Hello World").
// 				Expect().Status().Between(300, 400).
// 				Do()
// 		})
// 	})
// }
