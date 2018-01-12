FROM golang:latest
RUN mkdir /app
ADD . /app
ADD ./gosrc /go/src
ADD ./imports /go/src
WORKDIR /app
RUN go build -o main .
EXPOSE 80
ENTRYPOINT ./main 