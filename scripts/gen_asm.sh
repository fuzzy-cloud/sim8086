#!/bin/bash

set -e

for file in listings/*.asm; do
    base=$(basename "$file" .asm)
    nasm "$file" -o "listings/$base.bin"
    echo "Compiled $file -> listings/$base.bin"
done
