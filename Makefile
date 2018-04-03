HASH = $(shell git rev-parse --short HEAD)

cached: build 
	cd ./bin/ && ./gitdo -c

crosscompile:
	rm -rf ./out
	mkdir ./out
	env GOOS=windows GOARCH=386 vgo build -o ./out/gitdo_win_386.exe .
	env GOOS=windows GOARCH=amd64 vgo build -o ./out/gitdo_win_amd64.exe .
	env GOOS=darwin GOARCH=amd64 vgo build -o ./out/gitdo_mac_amd64 .
	env GOOS=darwin GOARCH=386 vgo build -o ./out/gitdo_mac_386 .

install: 
	vgo install -ldflags "-X main.version=0.0.1-$(HASH)"

test:
	vgo test github.com/nebloc/gitdo
	vgo test github.com/nebloc/gitdo/diffparse

release_dir:
	rm -rf release/
	mkdir release/
	env GOOS=windows GOARCH=386 vgo build -ldflags "-X main.version=0.0.1-$(HASH)" -o ./release/gitdo_win_i386.exe .
	env GOOS=windows GOARCH=amd64 vgo build -ldflags "-X main.version=0.0.1-$(HASH)" -o ./release/gitdo_win_amd64.exe .
	env GOOS=darwin GOARCH=amd64 vgo build -ldflags "-X main.version=0.0.1-$(HASH)" -o ./release/gitdo_darwin_amd64 .
	env GOOS=darwin GOARCH=386 vgo build -ldflags "-X main.version=0.0.1-$(HASH)" -o ./release/gitdo_darwin_i386 .
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp install.sh ./release/
	cp install.bat ./release/
