.PHONY: test fasttest clear generate

test: fasttest
	go test -v -count=1 -tags="doctest" doc_gen_test.go
	go test -v -count=1 -tags="doctest" readme_gen_test.go

fasttest: generate
	go test -v -count=1 ${GO_TEST_COVER_ARGS} ./...

clear:
	@-rm *_gen.go
	@-rm *_gen_test.go

generate: clear
	go run -v -tags="generate_numeric" ./generators/expect/numeric
	go run -v -tags="generate" ./generators/clear/clear
	go run -v -tags="generate" ./generators/clear/tests
	go run -v -tags="generate" ./generators/doc
	go run -v -tags="generate" ./generators/readme


