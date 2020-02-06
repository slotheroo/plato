FROM golang:1.10.3
WORKDIR /go/src/github.com/slotheroo/plato/platoserver
COPY . .
RUN go get -d -v ./...
RUN go build -v ./...
ENTRYPOINT ./platoserver
EXPOSE 3009
