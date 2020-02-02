# Basic go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Binary names
BINARY_NAME=shadowsocks-fyne.exe

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

RSRC=rsrc -manifest ./assets/app.manifest -o rsrc.syso

build-win64:
	$(GOGET) github.com/akavel/rsrc
	$(RSRC) -arch amd64
	#$(GOBUILD) -o $(BINARY_NAME) -v
	$(GOBUILD) -o $(BINARY_NAME) -ldflags -H=windowsgui -v

run-win64: build-win64
	FYNE_SCALE=2.5 ./$(BINARY_NAME)