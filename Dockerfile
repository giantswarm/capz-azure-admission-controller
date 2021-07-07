FROM alpine
WORKDIR /app
ADD build/manager build/manager
ENTRYPOINT ["manager"]
