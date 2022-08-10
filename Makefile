b:
	echo "Compiling..."
	go build -ldflags="-s -w" .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -o bin/main-freebsd .
	GOOS=linux GOARCH=amd64 go build -o bin/main-linux .
	GOOS=linux GOARCH=arm64 go build -o bin/main-linux .
	GOOS=darwin GOARCH=amd64 go build -o bin/main-darwin .
	GOOS=darwin GOARCH=arm64 go build -o bin/main-darwin-m1 .
