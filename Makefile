.ONESHELL:

lint:
	@command -v golangci-lint > /dev/null || go get github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run --exclude-use-default=false \
		--enable=golint \
		--enable=staticcheck \
		--enable=gosimple \
		--enable=unconvert \
		--enable=goconst \
		--enable=goimports \
		--enable=maligned \
		--enable=misspell \
		--enable=unparam \
		--enable=prealloc

test:
	@go test -race -v -failfast -cover ./...

bench:
	@go test -bench=. -benchmem

fuzz-clean:
	# Copy crashers to testdata for regression testing
	@find crashers -type f | grep -v '\.' | xargs -i cp {} testdata
	rm -rf corpus
	rm -rf crashers
	rm -rf suppressions

fuzz:
	@command -v go-fuzz || go get -u github.com/dvyukov/go-fuzz/...
	go-fuzz-build github.com/martingallagher/iter
	go-fuzz -bin=./iter-fuzz.zip -workdir=.
