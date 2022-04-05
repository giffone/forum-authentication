FROM golang

LABEL project="Forum"

WORKDIR /web

COPY . .

RUN go build cmd/forumsqlite/main.go

CMD ["./main"]