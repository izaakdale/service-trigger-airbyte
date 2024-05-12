FROM golang:1.22-alpine as builder
WORKDIR /
COPY . ./
RUN go mod download


RUN go build -o /service-trigger-airbyte


FROM alpine
COPY --from=builder /service-trigger-airbyte .


EXPOSE 80
CMD [ "/service-trigger-airbyte" ]
