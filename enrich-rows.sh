#!/bin/bash
f="${1:-}"
head -n1 "${f}"
address-enrich -input-file "${f}" -start-column="${2:-0}"
# head "${f}" |head |address-enrich -v
