BINARY	:= 3DC
MODULE	:= 3DC
VERSION	:= 0.5.0

GO		?= go
GOFLAGS	?=
LDFLAGS	?= -s -w

.PHONY: all build run test vet fmt tidy clean install uninstall help

all: build

build:
	$(GO) build $(GOFLAGS) -ldflags "$(LDFLAGS)" -o $(BINARY) .

run: build
	./$(BINARY) $(ARGS)

test:
	$(GO) test $(GOFLAGS) ./...

vet:
	$(GO) vet ./...

fmt:
	$(GO) fmt ./...

tidy:
	$(GO) mod tidy

clean:
	rm -f $(BINARY) $(BINARY).exe *.test coverage.out

install:
	$(GO) install $(GOFLAGS) -ldflags "$(LDFLAGS)" .

uninstall:
	rm -f "$$($(GO) env GOPATH)/bin/$(BINARY)"
	@if [ -n "$$($(GO) env GOBIN)" ]; then rm -f "$$($(GO) env GOBIN)/$(BINARY)"; fi

help:
	@echo "3D-Chess-CLI $(VERSION)"
	@echo ""
	@echo "Targets:"
	@echo "  build      Build ./$(BINARY) (default)"
	@echo "  run        Build and run (ARGS='game list')"
	@echo "  test       Run go test ./..."
	@echo "  vet        Run go vet ./..."
	@echo "  fmt        Run go fmt ./..."
	@echo "  tidy       Run go mod tidy"
	@echo "  clean      Remove binary and test artifacts"
	@echo "  install    Install to GOPATH/bin"
	@echo "  uninstall  Remove installed binary"
	@echo "  help       Show this message"
