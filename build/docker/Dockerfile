FROM scratch
MAINTAINER support@stocktopus.io
EXPOSE 443
COPY ca-certificates.crt /etc/ssl/certs/
COPY server /
CMD ["/server"]
