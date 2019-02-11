all:
	go build -v -o import.so -buildmode=c-shared .

clean:
	rm import.so import.h
