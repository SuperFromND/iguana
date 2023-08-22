all:
	go build -o bin/iguana.exe src/main.go

release:
	go build -o bin/iguana.exe -ldflags "-X main.version=1.0.0" src/main.go