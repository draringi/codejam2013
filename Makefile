.PHONEY: all deps kaguya

codejam2013: power_response.go deps
	go install draringi/codejam2013

all: codejam2013 deps

deps:
	go get github.com/fxsjy/RF.go/RF
	go install github.com/fxsjy/RF.go/RF
	go get github.com/jbarham/gopgsqldriver
	go install github.com/jbarham/gopgsqldriver

kaguya:
	cp -fr src ~/go/src/draringi/codejam2013
	cp -f power_response.go ~/go/src/draringi/codejam2013
	rm -f ~/go/bin/codejam2013
	rm -fr ~/go/pkg/linux_amd64/draringi
	go install draringi/codejam2013
	~/go/bin/codejam2013
