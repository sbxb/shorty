
.PHONY: run
run:
	go run cmd/shortener/main.go

.PHONY: test
test:
	go test -v -count=1 ./...
