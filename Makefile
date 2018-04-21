HASH = $(shell git rev-parse --short HEAD)
VERNUM=0.0.6-A5
VERSION=-ldflags "-X main.version=$(VERNUM)-$(HASH)"

install: 
	vgo build ${VERSION} -o /usr/local/bin/gitdo ./app/gitdo


test2:
	vgo test github.com/nebloc/gitdo/app/...

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_32.exe ./app/gitdo
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_64.exe ./app/gitdo
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_32 ./app/gitdo
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_64 ./app/gitdo
	env GOOS=linux GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_linux_32 ./app/gitdo
	env GOOS=linux GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_linux_64 ./app/gitdo
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp mac_linux_install.sh ./release/
	cp windows_install.bat ./release/
	cp ./docs/Install.md ./release/
	cp ./docs/Usage.md ./release/
