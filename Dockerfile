FROM ubuntu

ENV XDG_DATA_HOME=/home/.local/share
ENV XDG_STATE_HOME=/home/.local/state
RUN mkdir -p $XDG_DATA_HOME/polycloze

# Copy course files
COPY python/build/courses/eng-tgl.db $XDG_DATA_HOME/polycloze

# Copy executable
COPY polycloze /

ENTRYPOINT ["/polycloze"]
EXPOSE 3000
