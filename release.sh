#!/bin/bash
# -*- coding: utf-8 -*-

set -eux
: $1

git tag $1
GITHUB_TOKEN=`cat token.txt` goreleaser --rm-dist