FROM golang:1.16.3-alpine3.13

WORKDIR /go/src/app

RUN apk add --no-cache bind-tools
RUN go get github.com/githubnemo/CompileDaemon

ENV CGO_ENABLED 0

ENTRYPOINT REPLICA_ID="$(dig +short -x $(hostname -i))" \
    REPLICA_ID=${REPLICA_ID%%.*} REPLICA_ID=${REPLICA_ID##*_} \
    CompileDaemon --command=app --build="go build -o /go/bin/app"
