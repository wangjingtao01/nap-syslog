FROM hub.sky-cloud.net:5000/centos:7
LABEL maintainer "songtianyi <songtianyi@sky-cloud.net>"
COPY cmd/ekanited/ekanited /usr/bin/nap-syslog
CMD /usr/bin/nap-syslog -udp 0.0.0.0:5514 -tcp 0.0.0.0:5514 -dispatcher /etc/nap-syslog/dispatcher.json -queryhttp 0.0.0.0:8090
