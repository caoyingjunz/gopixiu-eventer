FROM golang:1.17.5 as builder
ARG GOPROXY
ARG APP
ENV GOPROXY=${GOPROXY}
WORKDIR /go/gopixiu-eventer
COPY . .
RUN CGO_ENABLED=0 go build -a -o ./dist/${APP} cmd/${APP}/${APP}.go

FROM jacky06/static:nonroot
ARG APP
WORKDIR /
COPY --from=builder /go/gopixiu-eventer/dist/${APP} /usr/local/bin/${APP}
USER root:root