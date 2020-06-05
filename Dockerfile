FROM alpine:3.6 as alpine
RUN apk add -U --no-cache ca-certificates

FROM scratch
MAINTAINER support@stocktopus.io
EXPOSE 443 80
COPY --from=alpine /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY ./bin/stocktopus /
CMD ["/stocktopus"]
