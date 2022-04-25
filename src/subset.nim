# Computes subset of relation (from stdin) containing items in file.

import std/os
import std/rdstdin
import std/sequtils
import std/sets
import std/streams
import std/strutils

proc main() =
  if paramCount() < 1:
    stderr.writeLine("missing argument: list of IDs")
    quit()

  var stream = openFileStream(paramStr(1))
  var ids = toHashSet(toSeq(stream.lines()))
  stream.close()

  var line: string
  while readLineFromStdin("", line) and line.len > 0:
    let row = line.split('\t')
    if row[0] in ids or row[1] in ids:
      echo line

main()
