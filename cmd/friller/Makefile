ROOT_DIR:=$(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

BINARY=friller
VERSION=1.0.0-BETA-6
SigFile=pronom:DROID_SignatureFile_V97.xml;container-signature-20201001.xml
BUILD=`git rev-parse --short HEAD`
PLATFORMS=darwin linux windows
ARCHITECTURES=amd64

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Build=${BUILD} -X main.SigFile=${SigFile} -w -s" 
BLDFLAGS=-buildmode=pie

all: clean build_all

build:
	go build ${BLDFLAGS} ${LDFLAGS} -o ${BINARY}

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(BLDFLAGS) $(LDFLAGS) -v -o $(BINARY)-$(GOOS))))
	mv friller-darwin friller && tar cvfz friller_macos_$(VERSION).tar.gz friller
	rm friller
	mv friller-linux friller && tar cvfz friller_linux_$(VERSION).tar.gz friller
	rm friller
	mv friller-windows friller && tar cvfz friller_win_$(VERSION).tar.gz friller
	rm friller


clean:
	rm -f '${BINARY}-linux'
	rm -f '${BINARY}-darwin'
	rm -f '${BINARY}-windows'

.PHONY: clean build build_all all
