FROM alpine
WORKDIR /app
ADD manager manager
ENTRYPOINT ["manager"]
