FROM golang:1.23.2 AS build
WORKDIR /go/src/github.com/natecw/minily
COPY . .

ENV CGO_ENABLED=0
RUN go build -o server .

FROM build as development
RUN apt-get update && apt-get install -y git
CMD ["go", "run", "main.go"]

FROM alpine:3.12
EXPOSE 8000
COPY --from=build /go/src/github.com/natecw/minily/server /server
CMD ["/server", "start", "--addr", ":8000"]