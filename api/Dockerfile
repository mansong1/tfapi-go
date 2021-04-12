FROM golang:alpine AS builder

# Create apiuser
ENV USER=api
ENV UID=10001

# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

RUN  mkdir -p /go/src \
  && mkdir -p /go/bin \
  && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/app
COPY . $GOPATH/src/app

# Fetch Certificates
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
        ca-certificates \
        git \
    && update-ca-certificates 2>/dev/null || true

WORKDIR $GOPATH/src/app
# Fetch dependencies
RUN go get github.com/gorilla/mux \
github.com/gorilla/handlers \
google.golang.org/grpc \
github.com/kelseyhightower/envconfig

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o api .

FROM scratch
EXPOSE 3000
# Import from builder
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/app/api /app/api
COPY --from=builder /go/src/app/conf/ /app/conf/

# Use an unprivileged user.
USER api:api

ENTRYPOINT ["/app/api"]