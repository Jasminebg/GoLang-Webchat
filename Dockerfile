# Stage 1
FROM golang:alpine as builder
RUN apk add bash ca-certificates git gcc g++ libc-dev
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app/backend 
RUN go mod download
RUN go build -v
RUN ls
CMD ["/app/Golang-Webchat"]
# FROM golang:alpine as builder
# RUN apk update && apk add --no-cache git
# RUN mkdir /build 
# ADD . /build/
# WORKDIR /build/backend
# RUN go get -d -v
# RUN go build -o backend .
# # Stage 2
# FROM alpine
# RUN adduser -S -D -H -h /app appuser
# USER appuser
# COPY --from=builder /build/backend/ /app/
# WORKDIR /app
# CMD ["./backend"]