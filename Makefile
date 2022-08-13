dev:
	echo "Compiling..."
	go build .
prod:
	echo "Compiling..."
	go build -ldflags="-s -w" .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-freebsd .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-linux .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/go-chaos-linux .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/go-chaos-darwin .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w"-o bin/go-chaos-darwin-m1 .
