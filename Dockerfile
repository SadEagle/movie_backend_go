FROM golang:tip-trixie AS builder
WORKDIR /backend_go
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN go build -o /server ./cmd/server/


FROM golang:tip-trixie
COPY --from=builder /server ./
# NOTE: temporal migration and user creation solution
# TODO: fix later
COPY ./script ./script
COPY ./db/migrations/ ./db/migrations

CMD [ "./server" ]

EXPOSE 8080
HEALTHCHECK CMD "wget --tries=1 --no-verbose http://backend:8080/healthcheck || exit 1"
