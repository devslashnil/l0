FROM golang:alpine as builder

RUN apk update && apk upgrade && \
    apk add --no-util git

RUN mkdir /app
WORKDIR /app

ENV GO111MODULE=on

COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pub .\\cmd\\sub\\main.go


FROM alpine:latest

RUN apk --no-util add ca-certificates

RUN mkdir /app
WORKDIR /app
ADD ./api/model.json /app/consignment.json
COPY --from=builder /app/pub .

CMD ["./pub"]