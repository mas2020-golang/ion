export GIT_COMMIT=$(shell git rev-list -1 --abbrev-commit HEAD)
export BUILD_DATE=$(shell date +%Y-%m-%d)
export INSTALL_PATH="/usr/local/bin"

test:
	@echo ">> ion test..."
	@go test github.com/mas2020-golang/ion/cmd/file -coverprofile=coverage.out

coverage:
	@go tool cover -html=coverage.out

goreleaser:
	@echo ">> start building..."
	@goreleaser  --rm-dist --snapshot --skip-publish
	@echo "done!"

install_on_mac: build test
	@echo ">> start install..."
	@echo ">> copying into ${INSTALL_PATH}..."
	@cp bin/ion-darwin-amd64 ${INSTALL_PATH}/ion
	@echo "done!"

run:
	clear
	go run main.go

build:
	@echo ">> compiling for every OS and Platform...${BUILD_DATE}"
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=${BUILD_DATE}" -o bin/ion-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=${BUILD_DATE}" -o bin/ion-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.BuildDate=${BUILD_DATE}" -o bin/ion-windows-amd64.exe main.go
	@echo "done!"

clean:
	@rm -rf bin
	@rm -rf dist

look_update_pkgs:
	@echo ">> take a look at the newer versions of dependency modules"
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null