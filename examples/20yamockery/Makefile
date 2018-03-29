default: mockery yamockery

mockery:
	mockery -dir ./test -name I

yamockery:
	go run main.go

setup:
	go get -v github.com/vektra/mockery/cmd/...
