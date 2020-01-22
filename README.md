# ORI Industries assignment

## Requirements
- git
- docker
- docker-compose
- golang:1.13
- dep
- mockgen
- protobuf-go

## Testing

Clone this repo and go to repo folder
```
git clone https://github.com/fedjo/ori-app.git && cd deevio-project
```

### Local

1. Compile proto files
```
cd pb
protoc  -I . --go_out=plugins=grpc:. ./*.proto
```
2. Download dependencies
```
cd ..
dep ensure -v
```
3. Set environment variables BIND\_PORT and SERVER\_ADDRESS
```
export BIND\_PORT=3000
export SERVER\_ADDRESS=localhost:3000
```
4. On a terminal run
```
go run -v srv/main.go
```
5. On another terminal run to calculate the sum
```
go run -v client/main.go sum 12 30
go run -v client/main.go sqrt 1764
go run -v client/main.go gcd 168 546
```


### Docker
1. Build docker image of the developed web service
```
cd docker
$ ./build.sh latest
```
2. Run docker-compose and examine logs
```
docker-compose up
```

## Documentation

Please read the [docs](doc/)

