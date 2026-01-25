FROM golang:tip-trixie AS go_build
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY ./. /backend_go
WORKDIR /backend_go
RUN go build -o backend_go


FROM golang:tip-trixie
COPY --from=go_build ./backend_go .
CMD [ "./backend_go" ]

EXPOSE 8080
HEALTHCHECK CMD "wget --tries=1 --no-verbose http://backend:8080/healthcheck || exit 1"
