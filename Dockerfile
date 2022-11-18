FROM ubuntu

ENV XDG_DATA_HOME=/home/.local/share
ENV XDG_STATE_HOME=/home/.local/state
RUN mkdir -p $XDG_DATA_HOME/polycloze/courses

# Copy some course files
COPY python/build/polycloze/courses/eng-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/spa-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-hrv.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-swe.db $XDG_DATA_HOME/polycloze/courses

# Copy executable
COPY polycloze /

ENTRYPOINT ["/polycloze"]
EXPOSE 3000
