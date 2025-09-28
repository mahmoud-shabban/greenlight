
build: ./server/api/main.go
	go build -o bin/app  ./server/api/

clean:
	rm -rf ./bin/*
