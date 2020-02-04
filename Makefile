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
	$(GOGET) fyne.io/fyne/cmd/fyne
	fyne bundle -package resources -prefix= ./resources > ./resources/resources_gen.go
	fyne package -os windows -name shadowsocks-fyne
