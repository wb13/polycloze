#!/usr/bin/env bash

# Strip everything but the sentences.

sed "s/.\+\t.\+\t\(.\+\)/\1/g"
