all: mac mac_amd64 linux windows

mac:
	GOOS=darwin GOARCH=arm64 go build -o builds/gwir ./
	tar -czvf builds/gwir.macos-arm64.tar.gz -C builds gwir
	rm -rf builds/gwir

mac_amd64:
	GOOS=darwin GOARCH=amd64 go build -o builds/gwir ./
	tar -czvf builds/gwir.macos-amd64.tar.gz -C builds gwir
	rm -rf builds/gwir

linux:
	GOOS=linux GOARCH=amd64 go build -o builds/gwir ./
	tar -czvf builds/gwir.linux-amd64.tar.gz -C builds gwir
	rm -rf builds/gwir

windows:
	GOOS=windows GOARCH=amd64 go build -o builds/gwir.exe ./
	tar -czvf builds/gwir.windows-amd64.tar.gz -C builds gwir.exe
	rm -rf builds/gwir.exe

clean:
	rm -rf builds/*
