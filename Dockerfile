FROM hub.sky-cloud.net:5000/centos:7
LABEL maintainer "songtianyi <songtianyi@sky-cloud.net>"
COPY cmd/ekanited/ekanited /usr/bin/nap-syslog
CMD /usr/bin/nap-syslog -udp localhost:5514 -tcp localhost:5514 -dispatcher /etc/nap-syslog/dispatcher.json -queryhttp localhost:8090
