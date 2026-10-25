# Makefile for A+ Engine
# make       -> build A+ Engine
# make check -> start test A+ Engine after compiling
# make run   -> run A+ Engine
# To build, need:
# + go
# + make

.PHONY: all run check

all:
	@echo "Building A+ Engine"
	go build -o a-engine .
	@echo "Saved to : a-engine. Run : ./a-engine. (Or make run)"

run:
	@echo "Starting A+ Engine"
	go run .

check:
	@echo "Starting test for A+ Engine"
	go test ./core/ -v
