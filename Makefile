.PHONY : dev prod install compile

dev:
	echo "Compiling..."
	go build -o go-chaos . 

prod:
	echo "Compiling..."
	go build -ldflags="-s -w" -o go-chaos .

move:
	mv go-chaos ~/bin/go-chaos

install: prod move

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-freebsd-amd64 .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/go-chaos-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/go-chaos-darwin-m1 .
