TARGET := ecm-sdk-go

SRCDIR := $(PWD)

CC := go
LFLAGS := -w -s

REPO := harbor.arfa.wise-paas.com/ensaasmesh/
TAG := 0.0.1.1

all: sdk-demo

sdk-demo:
	$(CC) build -mod=mod -ldflags '$(LFLAGS)' -o $(SRCDIR)/bin/demo $(SRCDIR)/example/main.go

sdk-demo-image:
	sudo docker build -t $(REPO)demo:$(TAG) -f $(SRCDIR)/example/Dockerfile .

.PHONY:clean
clean:
	rm -rf $(SRCDIR)/bin
