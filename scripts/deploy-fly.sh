#!/usr/bin/env bash

# Copyright (c) 2022 Levi Gruspe
# License: MIT, or AGPLv3 or later

flyctl auth login --verbose
flyctl deploy .
