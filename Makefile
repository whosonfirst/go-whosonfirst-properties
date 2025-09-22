GOMOD=$(shell test -f "go.work" && echo "readonly" || echo "vendor")
LDFLAGS=-s -w

cli:
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/report-properties cmd/report-properties/main.go
	go build -mod $(GOMOD) -ldflags="$(LDFLAGS)" -o bin/index-properties cmd/index-properties/main.go

docker:
	docker build -t whosonfirst-properties-indexing .	

