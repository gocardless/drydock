PREFIX=/usr/local
VERSION=0.0.4
BUILD_COMMAND=gom build -ldflags "-X main.version=$(VERSION)"

.PHONY: build clean

build:
	$(BUILD_COMMAND) -o drydock *.go

test:
	test -z "$(golint ./... | tee /dev/stderr)" || exit 1
	gom tool vet *.go || exit 1
	gom test -race -test.v . || exit 1

build-production: test
	GOOS=linux GOARCH=amd64 $(BUILD_COMMAND) -o drydock.linux_amd64 *.go

deb: build-production
	bundle exec fpm -s dir -t $@ -n drydock -v $(VERSION) \
		--description "Docker image cleaner" \
		--maintainer "GoCardless Engineering <engineering@gocardless.com>" \
		drydock.linux_amd64=$(PREFIX)/bin/drydock

clean:
	-rm -f drydock drydock.linux_amd64 cover.out drydock.test
	-rm -f drydock_${VERSION}_amd64.deb

