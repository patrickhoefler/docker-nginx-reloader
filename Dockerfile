FROM busybox
MAINTAINER patrick.hoefler@gmail.com
COPY docker-nginx-reloader /
ENTRYPOINT ["/docker-nginx-reloader"]
