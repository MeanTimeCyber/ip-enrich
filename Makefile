BINARY := ip-enrich
GO ?= go
LDFLAGS := -s -w
BUILD_FLAGS := -trimpath -ldflags="$(LDFLAGS)"
IP ?= 109.158.10.179

.PHONY: build run clean

build:
	$(GO) build $(BUILD_FLAGS) -o $(BINARY) ./cli

run:
	$(GO) run cli/*.go -i $(IP)

clean:
	rm -f $(BINARY)