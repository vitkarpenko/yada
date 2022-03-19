FROM golang:1.18-bullseye AS build
COPY . /app
WORKDIR /app
RUN go build .

FROM debian:bullseye-slim
RUN apt update && \
    apt install -y ca-certificates
COPY --from=build /app/yada /yada/binary
COPY --from=build /app/.env /yada/.env
WORKDIR /yada
CMD ["./binary"]