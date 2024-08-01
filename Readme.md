## Setup
For generating proto grpc
```bash
brew install protobuf
go install github.com/golang/protobuf/protoc-gen-go@latest
```

For generating swagger doc
```bash
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/swaggo/gin-swagger@latest
go install github.com/swaggo/files@latest
```
## Tasks
### Generate proto grpc
```bash
make generateProto
```

### Generate API gateway code
```bash
go run autogen/main.go
```
Remember to commit generated code as well.

### Generate swagger doc
```bash
swag init -g cmd/main.go
swag init --parseDependency -g cmd/main.go
```

### Generate key pairs
Generate JWT for signing and verification (testing on local)
```bash
# generate private key
openssl genrsa -out auth.key 4096
# generate public key
openssl rsa -in auth.key -outform PEM -pubout -out auth.pub
```
The `auth.key` can be used to generate JWT for testing.

Run `go run cmd/generateJwt.go` to generate a new JWT.

Put generated `auth.pub` file into `config/envs` folder and adapt environment configuration file accordingly.

## Running the app
```bash
# development
make server
```

Testing:
~~~
curl localhost:3000/api/login -d '{"email":"abc", "password":"pass"}'
curl -XPOST -H "Authorization: Bearer jwt" localhost:3000/product/create -d '{"name":"cake", "sku":"123", "stock":1, "price":2}'
curl -H "Authorization: Bearer testjwt12345678819" localhost:3000/product/10
~~~

~~~
go install -mod=mod github.com/githubnemo/CompileDaemon
go mod tidy
go mod verify
go mod vendor
CompileDaemon --build="go build cmd/main.go" --command=./main
~~~