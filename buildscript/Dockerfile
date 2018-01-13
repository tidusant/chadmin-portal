FROM scratch
#ADD ca-certificates.crt /etc/ssl/certs/
ADD main /
ADD config /config
CMD ["/main","--port=80","&"] 