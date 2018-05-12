FROM scratch
MAINTAINER support@stocktopus.io
EXPOSE 443
ADD ca-certificates.crt /etc/ssl/certs
ADD server /
CMD ["/server"]
