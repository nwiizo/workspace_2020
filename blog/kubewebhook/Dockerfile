FROM golang:1.14.3-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /bin/example .

FROM alpine:latest
COPY --from=build /bin/example /usr/local/bin/example
RUN apk --no-cache add ca-certificates
ENTRYPOINT ["/usr/local/bin/example"]
