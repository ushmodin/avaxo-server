all: build

build:
	mkdir -p dist
	cd cmd; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../dist/main

clean:
	rm ./main
