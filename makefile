VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-dev)?)' version.go)

.PHONY: build docker release

build:
	go fmt ./...
	CGO_ENABLED=0 go build -o bin/ach ./cmd/server

client:
# Download
	if [ ! -d "$(shell pwd)/tmp/swagger-codegen" ]; then \
		git clone https://github.com/swagger-api/swagger-codegen tmp/swagger-codegen; \
	fi
	cd tmp/swagger-codegen && \
	mvn clean package && \
	java -jar tmp/swagger-codegen/modules/swagger-codegen-cli/target/swagger-codegen-cli.jar \
	  generate -i server/openapi.yaml -l go -o client/

clean:
	@rm -rf tmp/

docker: clean
	docker build -t moov/ach:$(VERSION) -f Dockerfile .
	docker tag moov/ach:$(VERSION) moov/ach:latest

release: docker AUTHORS
	go test ./...
	git tag -f $(VERSION)

release-push:
	git push origin $(VERSION)
	docker push moov/ach:$(VERSION)

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@
