fmt:
	gofmt -w -s ./
	goimports -w  -local github.com/toysmars/distlib-go ./

gen-mock:
	mockery --all


test: gen-mock
	go test ./...
