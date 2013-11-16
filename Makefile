.PHONEY: all deps kaguya

web: web.go
	go build web.go

all: web deps

kaguya:
	export GOPATH=~/go
	export PORT=8080
	cp -R src ~/go/
	rm -f ~/go/bin/codejam2013
	rm -fr ~/go/pkg/linux_amd64/draringi
	go install draringi/codejam2013
	~/go/bin/codejam2013
