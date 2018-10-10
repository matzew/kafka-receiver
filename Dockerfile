FROM centos:7
MAINTAINER Matthias Wessendorf <matzew@apache.org>
ARG BINARY=./kreceiver

COPY ${BINARY} /opt/kreceiver
ENTRYPOINT ["/opt/kreceiver"]
