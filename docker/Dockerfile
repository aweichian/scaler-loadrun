FROM ic-harbor.baozun.com/ic/normal-base:v0.1

USER vmuser

RUN mkdir -p /opt/project/
ADD main /opt/project/main
RUN echo -e "sleep 20\n./main" > /opt/project/start.sh && chmod u+x /opt/project/start.sh
EXPOSE 8080

WORKDIR /opt/project

ENTRYPOINT ["sh","start.sh"]