HASH = $(shell git rev-parse --short HEAD)
VERNUM=0.0.0-A2
VERSION=-ldflags "-X main.version=$(VERNUM)-$(HASH)"

cached: build 
	cd ./bin/ && ./gitdo -c

install: 
	vgo install ${VERSION}

test:
	vgo test github.com/nebloc/gitdo
	vgo test github.com/nebloc/gitdo/diffparse

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_win_i386.exe .
	env GOOS=windows GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_win_amd64.exe .
	env GOOS=darwin GOARCH=amd64 vgo build ${VERSION} -o ./release/gitdo_darwin_amd64 .
	env GOOS=darwin GOARCH=386 vgo build ${VERSION} -o ./release/gitdo_darwin_i386 .
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp install.sh ./release/
	cp install.bat ./release/
