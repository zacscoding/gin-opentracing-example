FROM golang:1.13-alpine AS build

RUN mkdir -p /go/src/github.com/zacscoding/gin-opentracing ~/.ssh && \
    apk add --no-cache git openssh-client make gcc libc-dev
WORKDIR /go/src/github.com/zacscoding/gin-opentracing
COPY . .
RUN make build

FROM alpine:3
COPY --from=build /go/src/github.com/zacscoding/gin-opentracing/gin-opentracing /bin/gin-opentracing
CMD /bin/gin-opentracing
