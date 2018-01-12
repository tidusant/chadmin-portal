FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go get github.com/gin-gonic/gin
RUN go get github.com/spf13/viper

RUN echo "machine github.com\n\tlogin tidusread\n\tpassword Password2@" >> ~/.netrc

#RUN echo "[url \"git@github.com:\"]\n\tinsteadOf = https://github.com/" >> /root/.gitconfig
#RUN mkdir /root/.ssh && echo "StrictHostKeyChecking no " > /root/.ssh/config

#git clone https://703dce2b95f8e4775d968a94153ac5c8f4c81481:x-oauth-basic@github.com/tidusant/c3m-common /go/src/github.com/tidusant/c3m-common


RUN go get github.com/tidusant/c3m-common/...
RUN go get github.com/tidusant/chadmin-repo/...
RUN go build -o main .
EXPOSE 80
ENTRYPOINT ./main 