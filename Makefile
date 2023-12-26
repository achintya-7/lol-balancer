test_servers:
	cd test_servers && go run main.go

lb:
	go run main.go

test:
	curl localhost:2205
	