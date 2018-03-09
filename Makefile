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
