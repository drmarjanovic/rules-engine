FROM golang:1.8 AS builder
WORKDIR /go/src/github.com/MainfluxLabs/rules-engine
COPY . .
RUN cd cmd/ && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o exe

FROM scratch
COPY --from=builder /go/src/github.com/MainfluxLabs/rules-engine/cmd/exe /
EXPOSE 9000
ENTRYPOINT [ "/exe" ]
