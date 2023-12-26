# Lol Balancer
* It is a simple load balancer written in Go that can be used to distribute the load between multiple servers. It follows a simple round-robin algorithm to distribute the load after every request.

* It currently only supports passive health checks which run periodically to check if the server is up or not. If the server is down, it labels the server as Inactive and does not send any requests to it. If at any point the server is up, it labels it as Active and starts sending requests to it.

* The periocity of the health as of now is 20 seconds. {TODO: Make it configurable}

## How to run
To test it locally, you can run the following commands from the Makefile to build 2 dummy servers and the load balancer.
```bash
make test_server
make lb
```
The above commands will start 2 dummy servers on ports `3000` and `3001` and the load balancer on port `2205`. You can then send requests to the load balancer on port 8080 and it will distribute the load between the 2 servers.
To call the api, you can use the following command:
```bash
curl localhost:2205
```

### PS
As per the name, this is not a replacement for a production ready load balancer like NGINX. It is just a simple implementation to understand the concepts of load balancing and health checks.