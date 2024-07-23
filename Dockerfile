FROM --platform=$BUILDPLATFORM golang:1.21.6-alpine3.18 AS builder
ARG SCOPE
RUN apk add --no-cache git upx
WORKDIR /app
COPY ["go.mod","go.sum", "./"]
RUN go mod download -x
COPY . .
RUN if [ $SCOPE != "local" ] ; then export CONFIG_DIR=/app/pkg/config && export SCOPE=$SCOPE && go test -v ./...; fi
RUN go build api/main.go

#upload compilance
FROM alpine:3.18 AS runner
ARG SCOPE
RUN apk update
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/pkg/config/properties.yml .

ENTRYPOINT [ "./main" ]
# docker build -t furialfonso/cow_project:latest .
# docker build --build-arg SCOPE=local -t furialfonso/cow_project:latest .
# docker run -e PORT=9000 -p 8080:8080 furialfonso/cow_project:latest