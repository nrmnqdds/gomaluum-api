#!/bin/bash
export GOTRACEBACK=all

stdbuf -oL air -c .air.toml 2>&1 | \
stdbuf -oL panicparse -force-color -rel-path | \
stdbuf -oL sed 's/\(app\/[^ ]*\)/\1    /g' | \
stdbuf -oL sed 's/\(app\/\)//g'              

