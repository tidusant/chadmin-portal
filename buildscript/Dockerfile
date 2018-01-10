FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 
EXPOSE 80
RUN go build -o main . 
CMD ["/app/main"]
