# Dev environment
# Usage:
# docker build -t <image> .
# docker run -dit --name <container> <image>
# sudo docker exec <container> sh -c 'rm -rf "/src/*"'
# sudo docker cp . <container>:/src

FROM node:18
FROM golang:1.19

COPY --from=0 / /
WORKDIR /src

RUN apt-get update
RUN apt-get install shellcheck
RUN apt-get install sqlite3
