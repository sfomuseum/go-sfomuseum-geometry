cli:
	go build -mod vendor -o bin/union-featurecollection cmd/union-featurecollection/main.go
	go build -mod vendor -o bin/clip-featurecollection cmd/clip-featurecollection/main.go
	go build -mod vendor -o bin/difference-featurecollection cmd/difference-featurecollection/main.go
