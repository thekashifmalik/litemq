FROM golang:1.23
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY gen gen
COPY internal internal
COPY cmd cmd
RUN go build -o build/litemq cmd/litemq/main.go
CMD [ "build/litemq" ]
