.PHONEY: all deps

web: web.go
	go build web.go

all: web deps
