FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY demo-app/main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o demo-app-binary main.go

FROM alpine:latest

COPY --from=builder /build/demo-app-binary /bin/demo-app-binary

RUN chmod +x /bin/demo-app-binary

ENTRYPOINT ["/bin/demo-app-binary"]