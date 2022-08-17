# Usage:
# docker build -t <image> .
# docker run -dit --name <container> <image>
# docker exec <container> sh -c 'rm -rf "/src/*"'
# docker cp . <container>:/src

FROM python:3.10

RUN apt-get update
RUN apt-get install ripgrep
RUN apt-get install shellcheck
RUN apt-get install sqlite3
RUN pip install litecli

WORKDIR /src

COPY requirements requirements
RUN pip install -r requirements/dev.requirements.txt
