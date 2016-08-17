# Booking.com Back End CodeSprint

Solution for the open ended question from the [**Booking.com
Back End CodeSprint**](https://www.hackerrank.com/booking-passions-hacked-backend).

This is a proof of concept of a search engine devised for helping 
[Booking.com](http://booking.com) users to find destinations around the
world based on their passions. It's written in [golang](https://golang.org/)
and uses [gRPC](http://www.grpc.io/) and the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
to expose a high-efficiency RPC service and REST API.

A more detailed description of the approach used might be found inside
the **docs** folder.


## Usage

Once you have installed Go version 1.6 or later just run the commands:

```
$ go get -u github.com/paulocezar/booking-codesprint
$ booking-codesprint serve
```

<sub> **(If you're reading this before the submission deadline the code 
is not on github yet, so please copy the content of the zip file to 
`$GOPATH/src/github.com/paulocezar/booking-codesprint` and run 
`go install github.com/paulocezar/booking-codesprint` instead of the 
`go get -u ...`)** </sub>


And that's all. The service is now running via REST on port 1337 and RPC
on 13337. For more info run `booking-codesprint`. Also `booking-codesprint
serve -h` shows how the database being queried and the ports the service
run might be changed.

Once service is running Swagger definitions might be found at 
[http://localhost:1337/swagger/service.swagger.json](http://localhost:1337/swagger/service.swagger.json).
