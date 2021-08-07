FROM alpine:3.4

USER root

RUN mkdir -p /opt/project/
ADD scaler-loadrun /opt/project/scaler-loadrun
RUN echo -e "sleep 1\n./scaler-loadrun" > /opt/project/start.sh && chmod u+x /opt/project/start.sh
EXPOSE 8080

WORKDIR /opt/project

ENTRYPOINT ["sh","start.sh"]
