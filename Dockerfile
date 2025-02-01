# Building Backend
FROM golang:alpine as passport-server

WORKDIR /source
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -buildvcs -o /dist ./pkg/main.go

# Runtime
FROM golang:alpine

COPY --from=passport-server /dist /passport/server
COPY ./templates /templates
COPY ./locales /locales

EXPOSE 8444

CMD ["/passport/server"]
