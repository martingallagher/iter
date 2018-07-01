lint:
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install >/dev/null
	gometalinter ./... --vendor --skip=vendor --exclude=\.*_mock\.*\.go --exclude=vendor\.* --cyclo-over=20 --deadline=10m --disable-all \
	--enable=errcheck \
	--enable=vet \
	--enable=deadcode \
	--enable=gocyclo \
	--enable=golint \
	--enable=varcheck \
	--enable=structcheck \
	--enable=maligned \
	--enable=ineffassign \
	--enable=interfacer \
	--enable=unconvert \
	--enable=goconst \
	--enable=gosimple \
	--enable=staticcheck \
	--enable=gas

test:
	go test -v -failfast -cover

bench:
	go test -failfast -bench=. -benchmem

fuzz-clean:
	# Copy to testdata
	find crashers -type f | grep -v '\.' | xargs -i cp {} testdata
	rm -rf corpus
	rm -rf crashers
	rm -rf suppressions

fuzz:
	command -v go-fuzz || go get -u github.com/dvyukov/go-fuzz/...
	go-fuzz-build github.com/martingallagher/iter
	go-fuzz -bin=./iter-fuzz.zip -workdir=.
