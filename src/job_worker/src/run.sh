#!/bin/bash

# start weed mount
/usr/bin/weed mount -filer=filer:8888 -dir=/mnt/weed &&

# start python worker
/usr/bin/python src/worker.py