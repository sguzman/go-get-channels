FROM golang

RUN mkdir /app
ADD . /app/
WORKDIR /app

RUN go get -u -v "github.com/PuerkitoBio/goquery" \
    && go get -u -v "github.com/deckarep/golang-set" \
    && go get -u -v "github.com/lib/pq"

RUN go build -o main .
CMD ["/app/main"]