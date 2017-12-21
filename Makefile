.PHONY: run

SRC=main.go crimes pdf
MAIN_SRC=main.go

# Runs the server
run:
	go run ${MAIN_SRC}

# Adds all required go imports
imports:
	goimports -l -w ${SRC}

# Formats go source
fmt:
	gofmt -w ${SRC}
