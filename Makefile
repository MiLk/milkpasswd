.PHONY: all deps fmt vet test install

COMMANDS := $(wildcard commands/*)
DIRS:= . $(COMMANDS)

all: fmt vet test install

deps:
	go get github.com/boltdb/bolt
	go get github.com/smartystreets/goconvey/convey

fmt:
	@for dir in $(DIRS); do \
		echo "Fmt $$dir/*.go"; \
		go fmt $$dir/*.go; \
	done

vet:
	@for dir in $(DIRS); do \
		echo "Vet $$dir/*.go"; \
		go vet $$dir/*.go; \
	done

test:
	@go test

install:
	@for command in $(COMMANDS); do \
		echo "Install $$command"; \
	  cd $$command && go install; \
	done

