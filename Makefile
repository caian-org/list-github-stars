.DEFAULT_GOAL := build

ARTIFACT = lgs

ifeq ($(OS),Windows_NT)
	ARTIFACT = lgs.exe
endif


build:
	@go build -o $(ARTIFACT)

release:
	@CGO_ENABLED=0 \
		go build \
		-v \
		-trimpath \
		-gcflags=all="-l -B" \
		-ldflags="-w -s" \
		-o $(ARTIFACT)

run: build
	@./$(ARTIFACT)

format:
	gofmt -s -w .
