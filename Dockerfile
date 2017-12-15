FROM golang:1.8
COPY . /go/src/github.com/prydonius/mariadb-broker
WORKDIR /go/src/github.com/prydonius/mariadb-broker
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /mariadb-broker . && strip /mariadb-broker

FROM alpine:3.6
COPY --from=0 /mariadb-broker /mariadb-broker
CMD ["/mariadb-broker", "-logtostderr"]
