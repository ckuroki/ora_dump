all: ora_dump

ora_dump: ora_dump.go
	go build ora_dump.go
