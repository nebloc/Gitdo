HASH = $(shell git rev-parse --short HEAD)
DATE = $(shell date -u +.%Y%m%d.%H%M%S)
VERNUM=0.1.0
VERSION=-ldflags '-X "main.version=$(VERNUM) ($(HASH)$(DATE))"'

install: 
	vgo build ${VERSION} -o /usr/local/bin/gitdo .

build:
	vgo build ${VERSION} .
test:
	vgo test github.com/nebloc/gitdo/...

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_32.exe ./gitdo
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_64.exe ./gitdo
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_32 ./gitdo
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_64 ./gitdo
	env GOOS=linux GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_linux_32 ./gitdo
	env GOOS=linux GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_linux_64 ./gitdo
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp mac_linux_install.sh ./release/
	cp windows_install.bat ./release/
	cp ./docs/Install.md ./release/
	cp ./docs/Usage.md ./release/
