FROM golang:1.9.3
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/configmap-example
COPY . /go/src/configmap-example
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /configmap-example ./cmd/main.go

FROM scratch
COPY --from=0 /configmap-example .
ENTRYPOINT ["/configmap-example"]
