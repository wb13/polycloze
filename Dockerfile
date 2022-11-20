FROM ubuntu

ENV XDG_DATA_HOME=/home/.local/share
ENV XDG_STATE_HOME=/home/.local/state
RUN mkdir -p $XDG_DATA_HOME/polycloze/courses

# Copy some course files

COPY python/build/polycloze/courses/eng-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-hrv.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-tgl.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-tgl.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-epo.db $XDG_DATA_HOME/polycloze/courses
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
