FROM golang:1.13-alpine AS build-env
RUN apk --no-cache add build-base git bzr mercurial gcc
ENV D=/myapp
WORKDIR $D
# cache dependencies
ADD go.mod $D
ADD go.sum $D
RUN go mod download
# now build
ADD . $D
RUN cd $D && go build -o mybin && cp mybin /tmp/

# final stage
FROM alpine
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=build-env /tmp/mybin /app/myapp
CMD ["/app/myapp"]
