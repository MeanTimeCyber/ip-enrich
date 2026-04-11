BINARY := ip-enrich
GO ?= go
LDFLAGS := -s -w
BUILD_FLAGS := -trimpath -ldflags="$(LDFLAGS)"

.PHONY: build run clean

build:
	$(GO) build $(BUILD_FLAGS) -o $(BINARY) ./cli

clean:
	rm -f $(BINARY)