COMMANDS = agwlisten agwconnect

all: $(COMMANDS)

agwlisten: $(patsubst %, ./cmd/agwlisten/%, $(shell go list -f '{{ join .GoFiles " " }}' ./cmd/agwlisten/))
	go mod tidy
	go build ./cmd/agwlisten

agwconnect: $(patsubst %, ./cmd/agwconnect/%, $(shell go list -f '{{ join .GoFiles " " }}' ./cmd/agwconnect/))
	go mod tidy
	go build ./cmd/agwconnect

clean:
	rm -f $(COMMANDS)
