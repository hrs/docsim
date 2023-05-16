# Build parameters
GOCMD=go
GOBUILD=$(GOCMD) build
BINARY=docsim

# Install parameters
PREFIX=/usr/local
INSTDIR=$(DESTDIR)$(PREFIX)/bin
MANDIR=$(DESTDIR)$(PREFIX)/share/man/man1
MANPAGE=$(BINARY).1

.PHONY: build
build: $(BINARY)

$(BINARY): $(shell find . -iname *.go)
	$(GOBUILD) -o $(BINARY) -v ./...

.PHONY: install
install: $(BINARY)
	mkdir -p $(INSTDIR) $(MANDIR)
	cp $(BINARY) $(INSTDIR)
	install -m 644 man/$(MANPAGE) $(MANDIR)/$(MANPAGE)

.PHONY: uninstall
uninstall:
	rm -f $(INSTDIR)/$(BINARY)
	rm -f $(MANDIR)/$(MANPAGE)

.PHONY: deps
deps:
	go mod download

.PHONY: test
test:
	go test -race -v ./...
	go vet ./...

.PHONY: clean
clean:
	rm -f $(BINARY)
	go clean
