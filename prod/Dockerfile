FROM golang
COPY main.go main.go
RUN CGO_ENABLED=0 go build -o main main.go
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

FROM scratch
COPY /static /static
COPY --from=0 /go/main /main
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /etc_passwd /etc/passwd

USER nobody
CMD ["/main"]