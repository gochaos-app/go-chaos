dev:
	echo "Compiling..."
	go build -o chaosctl . 
prod:
	echo "Compiling..."
	go build -ldflags="-s -w" -o chaosctl .

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=amd64 go build -ldflags="-s -w" -o bin/chaosctl-freebsd .
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/chaosctl-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o bin/chaosctl-linux-arm64 .
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o bin/chaosctl-darwin .
	GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o bin/chaosctl-darwin-m1 .
