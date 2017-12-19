FROM scratch

COPY bin/ranch53 /usr/local/bin/ranch53

WORKDIR /usr/local/bin

ENTRYPOINT ["ranch53"]
