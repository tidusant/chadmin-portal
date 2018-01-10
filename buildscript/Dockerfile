FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
CMD ["ls"]
WORKDIR /app 
EXPOSE 80
CMD ["ls"]
RUN go build -o main ./app/ 
ENTRYPOINT ./app/main
