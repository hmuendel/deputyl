FROM golang:latest as builder
COPY . /go/src/gitlab.com/hmuendel/deputyl
WORKDIR /go/src/gitlab.com/hmuendel/deputyl
RUN CGO_ENABLED=0 GOOS=linux go build --ldflags "-X main.VERSION=$(git describe) -X main.COMMIT=$(git rev-parse HEAD) -s -w" -a -installsuffix cgo .

FROM alpine
COPY --from=builder /go/src/gitlab.com/hmuendel/deputyl/deputyl deputyl
EXPOSE 8080
ENTRYPOINT ["/deputyl"]




