FROM alpine
COPY main /
VOLUME /config.json
EXPOSE 5000-15000
CMD exec /main
