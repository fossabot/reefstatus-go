FROM golang

COPY .  /go/src/github.com/cjburchell/reefstatus-go

WORKDIR  /go/src/github.com/cjburchell/reefstatus-go

RUN go build -o main .

CMD ["/go/src/github.com/cjburchell/reefstatus-go/main"]
