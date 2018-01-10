FROM golang:latest 
RUN mkdir /app 
ADD . /app
WORKDIR /app 
EXPOSE 80
CMD ["ls"]
RUN go build -o main .
ENTRYPOINT ./main
