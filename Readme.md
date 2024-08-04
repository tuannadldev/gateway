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

### Generate proto grpc
```bash
make protos
```

### Generate API gateway code for handle restful request and call grpc server
```bash
go run autogen/main.go
```
### Generate swagger doc
```bash
swag init -g cmd/main.go
swag init --parseDependency -g cmd/main.go
```

### Generate key pairs
Generate JWT for signing and verification (testing on local)
```bash
# generate private key
openssl genrsa -out cert.key 4096
# generate public key
openssl rsa -in cert.key -outform PEM -pubout -out cert.pub
```
The `cert.key` can be used to generate JWT.
Put generated `cert.pub` file into `config/envs` folder and adapt environment configuration file accordingly.

## Running the app
```bash
# development
make server
```