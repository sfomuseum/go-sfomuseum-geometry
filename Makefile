GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/union-featurecollection cmd/union-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/difference-featurecollection cmd/difference-featurecollection/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/union-features cmd/union-features/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/difference-features cmd/difference-features/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/clip-featurecollection cmd/clip-featurecollection/main.go

