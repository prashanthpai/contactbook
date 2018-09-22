.PHONY: all build clean

BUILDDIR = build
SERVER_BIN = contactbook

all: build

build:
	dep ensure --vendor-only
	go build -o ${BUILDDIR}/${SERVER_BIN} .

install: build
	install -D ${BUILDDIR}/${SERVER_BIN} ${GOPATH}/bin

clean:
	rm -f ${BUILDDIR}/${SERVER_BIN}

uninstall:
	rm -f ${GOPATH}/bin/${SERVER_BIN}

test: build
	go test -race -v ./...

functest: build
	go test -race -v -tags=integration ./...
