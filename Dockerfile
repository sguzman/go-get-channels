FROM golang as base

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get -u "github.com/PuerkitoBio/goquery"
RUN go get -u "github.com/deckarep/golang-set"
RUN go get -u "github.com/lib/pq"
RUN go build -o main .

FROM alpine

COPY --from=base /app/main /main
ENTRYPOINT ["/bin/sh"]