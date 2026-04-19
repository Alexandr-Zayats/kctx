build:
	go build -o bin/kctx

install:
	cp bin/kctx /usr/local/bin/kctx

run:
	go run main.go switch
