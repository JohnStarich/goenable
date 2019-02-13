all:
	go build -v -o namespace.so -buildmode=c-shared .

clean:
	rm namespace.so namespace.h
