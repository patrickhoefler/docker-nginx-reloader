FROM busybox
MAINTAINER patrick.hoefler@gmail.com
COPY bin/linux-amd64/docker-nginx-reloader /
CMD ["/docker-nginx-reloader"]
