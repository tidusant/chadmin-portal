FROM golang:latest
RUN mkdir /app
ADD ./chadmin-portal /app
ADD ./gosrc /go/src
ADD ./imports /go/src
WORKDIR /app
RUN go build -o portal .
RUN rm -R /app
RUN rm -R /go/src
EXPOSE 80
ENTRYPOINT portal