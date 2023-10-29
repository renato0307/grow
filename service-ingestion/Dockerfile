FROM golang:1.21 as build

WORKDIR /go/src/svc
COPY . .

RUN go mod tidy
RUN go vet -v
RUN go test -v

RUN CGO_ENABLED=0 go build -o /go/bin/svc

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/svc /
ENTRYPOINT ["/svc"]