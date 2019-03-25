default: binary

binary: build

build_linux_x64:
			 GO111MODULE=on GOARCH=amd64 GOOS=linux go build -o dist/bin/eventreplay-linux-x64

build_darwin_x64:
			 GO111MODULE=on GOARCH=amd64 GOOS=darwin go build -o dist/bin/eventreplay-macos-x64

build_windows_x64:
			 GO111MODULE=on GOARCH=amd64 GOOS=windows go build -o dist/bin/eventreplay-windows-x64.exe

build: dist build_linux_x64	build_darwin_x64 build_windows_x64 copy_assets

copy_assets:
			 cp plugin.json dist/
			 cp -r samples dist/

clean:
			 rm -rf dist

dist:
			 mkdir -p dist/bin

.PHONY: binary build build_linux_x64 build_darwin_x64 build_windows_x64 clean dist