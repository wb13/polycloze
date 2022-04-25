# Simplifies reflexive relation.

import std/rdstdin
import std/strformat
import std/strutils

proc main() =
  var line: string
  while readLineFromStdin("", line) and line.len > 0:
    let row = line.split('\t')
    let source = parseInt(row[0])
    let target = parseInt(row[1])

    if source < target:
      writeLine(stdout, &"{source}\t{target}")
  flushFile(stdout)

main()
