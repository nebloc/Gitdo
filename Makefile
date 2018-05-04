HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u +.%Y%m%d.%H%M%S)
VERNUM=0.1.0
VERSION=-ldflags '-X "main.version=$(VERNUM) ($(HASH)$(DATE))"'

install: 
	cp -r ./resources/hooks ~/.gitdo/
	cp -r ./resources/plugins ~/.gitdo/
	vgo build ${VERSION} -o /usr/local/bin/gitdo .

build:
	vgo build ${VERSION} .
test:
	vgo test github.com/nebloc/gitdo/...

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_32.exe .
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_64.exe .
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_32 .
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_64 .
	env GOOS=linux GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_linux_32 .
	env GOOS=linux GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_linux_64 .
	cp -r ./resources/hooks ./release/
	cp -r ./resources/plugins ./release/
	cp -r ./resources/install_scripts/* ./release/
	cp ./resources/docs/Install.md ./release/
	cp ./resources/docs/Usage.md ./release/
