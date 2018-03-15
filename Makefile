init: 
	mkdir bin
	mkdir pkg
clean:
	rm -rf ./bin

build: clean
	mkdir bin
	go build -o ./bin/gitdo ./src/
	cp ./src/config.json ./bin/
	
run: build
	cd ./bin/ && ./gitdo

cached: build 
	cd ./bin/ && ./gitdo -c

crosscompile:
	rm -rf ./out
	mkdir ./out
	env GOOS=windows GOARCH=386 vgo build -o ./out/gitdo_win_386.exe .
	env GOOS=windows GOARCH=amd64 vgo build -o ./out/gitdo_win_amd64.exe .
	env GOOS=darwin GOARCH=amd64 vgo build -o ./out/gitdo_mac_amd64 .
	env GOOS=darwin GOARCH=386 vgo build -o ./out/gitdo_mac_386 .
