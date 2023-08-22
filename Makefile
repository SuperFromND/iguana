all:
	go build -trimpath -o bin/iguana.exe src/main.go

release:
	go build -trimpath -o bin/iguana.exe -ldflags "-X main.version=1.0.1" src/main.go