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
	./install.sh

test:
	vgo test github.com/nebbers1111/gitdo
	vgo test github.com/nebbers1111/gitdo/diffparse

release:
	mkdir release/
	vgo build -o release/gitdo .
	cp -r ./hooks ./release/
	cp -r ./plugins ./release/
	cp install.sh ./release/
	echo "{"trello_key":"","trello_token":""}" > ./release/secrets.json
