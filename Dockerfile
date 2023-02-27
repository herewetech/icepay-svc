# Static go service
FROM alpine:latest

ADD bin/* /opt/icepay/
WORKDIR /opt/icepay
EXPOSE 9900
CMD [ "/opt/icepay/svc" ]
