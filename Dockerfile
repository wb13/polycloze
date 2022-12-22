FROM ubuntu

ENV XDG_DATA_HOME=/home/.local/share
ENV XDG_STATE_HOME=/home/.local/state
RUN mkdir -p $XDG_DATA_HOME/polycloze/courses

# Copy some course files
COPY python/build/polycloze/courses/cat-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-ell.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-hrv.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-mkd.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-rus.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-tgl.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-tok.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/cat-ukr.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-cat.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ell.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-hrv.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-mkd.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-rus.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-tgl.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-tok.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/eng-ukr.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-cat.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ell.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-mkd.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-rus.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-tgl.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-tok.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/hrv-ukr.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-cat.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-dan.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-deu.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ell.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-eng.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-epo.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-fin.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-fra.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-hrv.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ita.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-lit.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-mkd.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-nld.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-nob.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-pol.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-por.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ron.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-rus.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-spa.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-swe.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-tok.db $XDG_DATA_HOME/polycloze/courses
COPY python/build/polycloze/courses/tgl-ukr.db $XDG_DATA_HOME/polycloze/courses

# Copy version file
COPY python/build/polycloze/version.txt $XDG_DATA_HOME/polycloze

# Copy executable
COPY polycloze /

ENTRYPOINT ["/polycloze"]
EXPOSE 3000
