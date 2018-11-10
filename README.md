# SAI
A Simple proxy server to generate specs of APIs 

### Development
first of all build the Image
```sh
$ make build
```
Then start a container
```sh
$ make run
```
Then up and run the server from inside the container
```sh
$ go run main.go
```
Also you get into a running container by:
```sh
$ make shell
```