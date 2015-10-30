
Download:

	go get github.com/Sogilis/go-todos

Build:

	go build .
	./go-todos

Commants to test:

	curl -v -X POST -d "premier todo" http://localhost:8080/list
	curl -v -X POST -d "2e todo" http://localhost:8080/list
	curl -v http://localhost:8080/list
	curl -v -X DELETE  http://localhost:8080/list/0
	curl -v http://localhost:8080/list
	curl -v -X POST -d "premier todo" http://localhost:8080/list
	curl -v -X POST -d "2e todo" http://localhost:8080/list
	curl -v http://localhost:8080/list

