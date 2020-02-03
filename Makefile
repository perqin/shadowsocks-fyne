# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=shadowsocks-fyne

build: linux win64

linux:
	$(GOBUILD) -o $(BINARY_NAME)-$@ -v

win64:
	$(GOGET) github.com/akavel/rsrc
	rsrc -manifest ./assets/app.manifest -o rsrc.syso -arch amd64
	$(GOBUILD) -o $(BINARY_NAME)-$@.exe -ldflags -H=windowsgui -v
