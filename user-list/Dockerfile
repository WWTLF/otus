# build stage
FROM golang as builder

ENV GO111MODULE=on

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ./cmd/user-list-server 

# final stage
FROM scratch
COPY --from=builder /app/user-list-server /app/
EXPOSE 8000
CMD ["/app/user-list-server", "--port=8000", "--host=0.0.0.0"]