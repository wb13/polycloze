FROM ubuntu

ENV XDG_DATA_HOME=/home/.local/share
ENV XDG_STATE_HOME=/home/.local/state
RUN mkdir -p $XDG_DATA_HOME/polycloze

# Copy some course files
COPY python/build/courses/eng-deu.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/eng-fra.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/eng-spa.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/spa-eng.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-dan.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-deu.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-eng.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-fin.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-fra.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-hrv.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-ita.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-lit.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-nld.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-nob.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-pol.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-por.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-ron.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-spa.db $XDG_DATA_HOME/polycloze
COPY python/build/courses/tgl-swe.db $XDG_DATA_HOME/polycloze

# Copy executable
COPY polycloze /

ENTRYPOINT ["/polycloze"]
EXPOSE 3000
