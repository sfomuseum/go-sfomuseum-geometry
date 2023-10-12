GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")

cli:
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/union-featurecollection cmd/union-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/clip-featurecollection cmd/clip-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="-s -w" -o bin/difference-featurecollection cmd/difference-featurecollection/main.go
