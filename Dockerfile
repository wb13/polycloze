# Dev environment
# Usage:
# docker build -t <image> .
# docker run -dit --name <container> <image>
# docker exec <container> sh -c 'rm -rf "/src/*"'
# docker cp . <container>:/src

FROM node:18
FROM golang:1.19

COPY --from=0 / /
WORKDIR /src

RUN apt-get update
RUN apt-get install shellcheck
RUN apt-get install sqlite3

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY api/js/package.json .
COPY api/js/package-lock.json .
RUN npm ci
