# build step
FROM golang:1.21.3 AS build

WORKDIR /app

COPY . /app

RUN CGO_ENABLED=0 GOOS=linux go build -o api main.go

# image step
FROM scratch

WORKDIR /app

COPY --from=build /app/api /app

EXPOSE 8000

CMD [ "./api" ]