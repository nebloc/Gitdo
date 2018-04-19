HASH = $(shell git rev-parse --short HEAD)
VERNUM=0.0.0-A5
VERSION=-ldflags "-X main.version=$(VERNUM)-$(HASH)"

cached: build 
	cd ./bin/ && ./gitdo -c

install: 
	vgo build ${VERSION} .
	mv gitdo /usr/local/bin/

test:
	vgo test github.com/nebloc/gitdo
	vgo test github.com/nebloc/gitdo/diffparse

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_32.exe .
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_64.exe .
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_32 .
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_64 .
	env GOOS=linux GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_linux_32 .
	env GOOS=linux GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_linux_64 .
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp mac_linux_install.sh ./release/
	cp windows_install.bat ./release/
	cp ./docs/Install.md ./release/
	cp ./docs/Usage.md ./release/
