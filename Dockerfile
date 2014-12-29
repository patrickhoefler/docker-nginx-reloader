FROM busybox
MAINTAINER patrick.hoefler@gmail.com
COPY docker-nginx-reloader /
CMD ["/docker-nginx-reloader"]
