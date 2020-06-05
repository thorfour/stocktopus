FROM scratch
MAINTAINER support@stocktopus.io
EXPOSE 443 80
COPY ca-certificates.crt /etc/ssl/certs/
COPY ./bin/stocktopus /
CMD ["/stocktopus"]
