#!/usr/bin/env bash

set -e

cd tools
perl -ne "print \$1.\"\n\" if /\"([a-zA-Z0-9\.\/_-]+$1)\"/" < tools.go | xargs -I {} go install {}
