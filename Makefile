COMMANDS = agwlisten agwconnect agwtalk

%: cmd/%/main.go
	go mod tidy
	go build ./cmd/$@

all: $(COMMANDS)

agwlisten: $(patsubst %, ./cmd/agwlisten/%, $(shell go list -f '{{ join .GoFiles " " }}' ./cmd/agwlisten/))

agwconnect: $(patsubst %, ./cmd/agwconnect/%, $(shell go list -f '{{ join .GoFiles " " }}' ./cmd/agwconnect/))

clean:
	rm -f $(COMMANDS)
