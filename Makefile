build:
	CGO_ENABLED=0 go build -o booking ./cmd/booking/

run: build
	./booking

test:
	go test ./...

clean:
	rm -f booking

.PHONY: build run test clean
