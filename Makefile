build: 
	go build -o bin/cache

run: build
	./bin/cache

runa: build
	./bin/cache --listenaddr :4000 --leaderaddr :3000