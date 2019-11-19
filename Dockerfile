FROM golang:1.13 as build-env
RUN mkdir -p /var/local/apirate/keys
RUN chmod -R 0700 /var/local/apirate
WORKDIR /src
COPY src .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o apirate .
RUN strip --strip-unneeded apirate

FROM scratch
COPY --from=build-env /src/apirate /
COPY --from=build-env /var/local/apirate /var/local/apirate
EXPOSE 8080
CMD ["./apirate"]
