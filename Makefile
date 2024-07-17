EXECNAME = iguana

ifeq ($(OS),Windows_NT)
	EXECNAME := $(addsuffix .exe,$(EXECNAME))
endif

all:
	go build -trimpath -o bin/$(EXECNAME) src/main.go

release:
	go build -trimpath -o bin/$(EXECNAME) -ldflags "-X main.version=1.3.1" src/main.go