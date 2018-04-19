HASH = $(shell git rev-parse --short HEAD)
VERNUM=0.0.0-A5
VERSION=-ldflags "-X main.version=$(VERNUM)-$(HASH)"

install: 
	vgo build ${VERSION} -o /usr/local/bin/gitdo ./app

test:
	vgo test github.com/nebloc/gitdo/app
	vgo test github.com/nebloc/gitdo/app/diffparse
	vgo test github.com/nebloc/gitdo/app/utils
	vgo test github.com/nebloc/gitdo/app/versioncontrol

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_32.exe ./app
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_64.exe ./app
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_32 ./app
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_64 ./app
	env GOOS=linux GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_linux_32 ./app
	env GOOS=linux GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_linux_64 ./app
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp mac_linux_install.sh ./release/
	cp windows_install.bat ./release/
	cp ./docs/Install.md ./release/
	cp ./docs/Usage.md ./release/
