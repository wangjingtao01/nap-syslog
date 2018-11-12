FROM alpine:3.8

LABEL maintainer "songtianyi <songtianyi@sky-cloud.net>"

RUN apk update && apk add tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    mkdir -p /etc/nap-syslog/ 

COPY cmd/ekanited/ekanited /usr/bin/nap-syslog
COPY dispatcher.json /etc/nap-syslog/

CMD /usr/bin/nap-syslog -udp 0.0.0.0:5514 -tcp 0.0.0.0:5514 -dispatcher /etc/nap-syslog/dispatcher.json -queryhttp 0.0.0.0:8090 -input rfc3164 -batchsize 2 -batchtime 1000
