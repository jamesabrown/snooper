FROM golang:1.9 AS build
LABEL maintainer "Peter Benjamin"
WORKDIR /go/src/github.com/jamesabrown/snooper
COPY . .
RUN go get -u -v github.com/golang/dep/cmd/dep \
    && dep init \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o snooper .

FROM alpine
COPY --from=build /go/src/github.com/jamesabrown/snooper/snooper /usr/bin/snooper
ENTRYPOINT [ "/usr/bin/snooper" ]
