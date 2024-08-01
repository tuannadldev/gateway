# build stage
FROM golang:1.22 as builder
EXPOSE 3000
ENV GO111MODULE=on
RUN mkdir -p /go/src/gateway
WORKDIR /go/src/gateway


COPY . .

# RUN apt-get update -y && apt-get upgrade -y
# RUN apt install curl && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
RUN go mod init
RUN go mod tidy
RUN go mod verify
RUN go mod vendor

ENV ELASTIC_APM_LOG_FILE=stderr
ENV ELASTIC_APM_LOG_LEVEL=debug
ENV ELASTIC_APM_SERVICE_NAME=gateway-service
ENV ELASTIC_APM_SERVER_URL=http://apm-server.apm.svc.service-name.dev

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build /go/src/gateway/cmd/main.go
## run dev

CMD ["./main"]

# # final stage => run production
# FROM scratch
# COPY --from=builder /go/src/gateway .
# EXPOSE 3000
# ENTRYPOINT ["./main"]