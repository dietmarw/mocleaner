# This Makefile is mainly useful for cross-compilation
# and upload to the GitHub release sections
# It is based on https://github.com/aktau/github-release/Makefile
# In order to be able to cross-compile you need to have
# built GO for all the compilation targets.
# Under Ubuntu Linux those are already available in the repo:
#  $ sudo apt-get install golang-$GOOS-$GOARCH
# For upload you need to have github-release installed
#  $ go get github.com/aktau/github-release

LAST_TAG := $(shell git describe --abbrev=0 --tags)

USER := dietmarw
EXECUTABLE := mocleaner
GITHUB_RELEASE = github-release

# only include the amd64 binaries, otherwise the github release will become
# too big
UNIX_EXECUTABLES := \
	darwin/amd64/$(EXECUTABLE) \
	freebsd/386/$(EXECUTABLE) \
	freebsd/amd64/$(EXECUTABLE) \
	linux/386/$(EXECUTABLE) \
	linux/amd64/$(EXECUTABLE)
WIN_EXECUTABLES := \
	windows/amd64/$(EXECUTABLE).exe

COMPRESSED_EXECUTABLES=$(UNIX_EXECUTABLES:%=%.tar.xz) $(WIN_EXECUTABLES:%.exe=%.zip)
COMPRESSED_EXECUTABLE_TARGETS=$(COMPRESSED_EXECUTABLES:%=build/%)

UPLOAD_CMD = $(GITHUB_RELEASE) upload -u $(USER) -r $(EXECUTABLE) -t $(LAST_TAG) -n $(subst /,-,$(FILE)) -f build/$(FILE)

all: $(EXECUTABLE)

# 386
build/darwin/386/$(EXECUTABLE):
	GOARCH=386 GOOS=darwin go build -o "$@"
build/freebsd/386/$(EXECUTABLE):
	GOARCH=386 GOOS=freebsd go build -o "$@"
build/linux/386/$(EXECUTABLE):
	GOARCH=386 GOOS=linux go build -o "$@"
build/windows/386/$(EXECUTABLE):
	GOARCH=386 GOOS=windows go build -o "$@"

# amd64
build/freebsd/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=freebsd go build -o "$@"
build/darwin/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=darwin go build -o "$@"
build/linux/amd64/$(EXECUTABLE):
	GOARCH=amd64 GOOS=linux go build -o "$@"
build/windows/amd64/$(EXECUTABLE).exe:
	GOARCH=amd64 GOOS=windows go build -o "$@"

# compressed artifacts, makes a huge difference (Go executable is ~9MB,
# after compressing ~2MB)
%.tar.xz: %
	tar -Jcvf "$<.tar.xz" "$<"
%.zip: %.exe
	zip -j "$@" "$<"

# git tag -a v$(RELEASE) -m 'release $(RELEASE)'
release: $(EXECUTABLE) $(COMPRESSED_EXECUTABLE_TARGETS)
	git push && git push --tags
	 $(GITHUB_RELEASE) release -u $(USER) -r $(EXECUTABLE) \
		-t $(LAST_TAG) -n $(LAST_TAG) || true
	$(foreach FILE,$(COMPRESSED_EXECUTABLES),$(UPLOAD_CMD);)

# install and/or update all dependencies, run this from the project directory
# go get -u ./...
# go test -i ./
dep:
	go list -f '{{join .Deps "\n"}}' | xargs go list -e -f '{{if not .Standard}}{{.ImportPath}}{{end}}' | xargs go get -u

$(EXECUTABLE): dep
	go build -o "$@"

install:
	go install

clean:
#	rm go-app || true
	rm $(EXECUTABLE) || true
	rm -rf build/

.PHONY: clean release dep install
