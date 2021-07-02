#build stage
FROM golang:alpine3.13 AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
# RUN go get -d -v ./...
RUN go build -o /go/bin/app -v ./...

#final stage
FROM alpine:3.13
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=jd Version=0.0.1
EXPOSE 4000
CMD [ "JD_central" ]
