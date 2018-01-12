FROM golang:latest
RUN mkdir /app
ADD ./chadmin-repo /app
ADD ./gosrc /go/src
ADD ./imports /go/src
WORKDIR /app
RUN go build -o portal .
EXPOSE 80
ENTRYPOINT portal