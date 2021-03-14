#!/bin/bash
f="${1:-}"
head -n1 "${f}"
addressenrich -input-file "${f}" -start-column="${2:-0}"
# head "${f}" |head |addressenrich -v 
