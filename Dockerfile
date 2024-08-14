FROM golang:bullseye AS build
WORKDIR /boardgames
ADD . .
RUN go build -o server main.go

FROM golang:bullseye
RUN apt-get update && apt-get install -y --no-install-recommends ffmpeg wget ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /boardgames/server /usr/bin/program
ENTRYPOINT ["/usr/bin/program"]