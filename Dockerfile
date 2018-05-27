FROM golang as build
COPY . ${GOPATH}/src/github.com/ushmodin/avaxo-server
WORKDIR ${GOPATH}/src/github.com/ushmodin/avaxo-server
RUN make 


FROM alpine
COPY --from=build /go/src/github.com/ushmodin/avaxo-server/dist/main /main
VOLUME /config.json
EXPOSE 5000-15000
CMD exec /main
