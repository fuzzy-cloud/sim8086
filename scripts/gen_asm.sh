#!/bin/bash

set -e

for file in listings/*.asm; do
    base=$(basename "$file" .asm)
    nasm "$file"
    echo "Compiled $file -> listings/$base"
done
