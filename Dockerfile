FROM alpine
WORKDIR /app
ADD build/manager ./
ENTRYPOINT ["manager"]
