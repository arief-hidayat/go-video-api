FROM golang:1.18 AS builder
WORKDIR /app
COPY go.mod .
COPY go.sum .
RUN go env -w GOPROXY=direct && go mod download
COPY controllers/ ./controllers/
COPY models/ ./models/  
COPY query/ ./query/  
COPY server/ ./server/  
COPY assets/ ./assets/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/api  ./server/cmd/main.go

FROM scratch
COPY --from=builder /go/bin/api /go/bin/api
COPY --from=builder /app/assets ./assets/
EXPOSE 8000
ENTRYPOINT ["/go/bin/api"]