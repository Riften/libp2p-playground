.PHONY:p2p
.PHONY:p2p.exe
,PHONY:p2p-win
init:
	mkdir bin
p2p:
	go build  -ldflags "-s -w" -o bin/p2p github.com/Riften/libp2p-playground
p2p.exe:
	go build  -ldflags "-s -w" -o bin/p2p.exe github.com/Riften/libp2p-playground
p2p-win:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build  -ldflags "-linkmode external -extldflags -static -s -w" -o bin/p2p.exe github.com/Riften/libp2p-playground