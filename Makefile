all: install test

TARGETS = echo cat

test:
	go test -v ./...

install:
	@for cmd in $(TARGETS); do \
	rm $(GOPATH)/bin/$$cmd; \
	done; \
	go install ./cmd/...
	@for cmd in $(TARGETS); do \
		echo -ne "\t\t"; \
		du -h $(GOPATH)/bin/$$cmd ; \
		echo -ne "stripped:\t"; \
		strip -s $(GOPATH)/bin/$$cmd ; \
		du -h $(GOPATH)/bin/$$cmd ; \
	done
