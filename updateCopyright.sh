#!/bin/bash

# Script to update the copyright line in all files to current year and version
#Mostly Ai generated

directories=(
    "cmd"
    "config"
    "data"
    "internal"
    "util"
    "webView"
)

copyright="// Copyright © 2026 Noah Yurasko distributed under GNU GENERAL PUBLIC LICENSE V3"

for rootDir in "${directories[@]}"; do
    while IFS= read -r file; do
        echo "Checking $file"

        # Check if the file already has a copyright header
        if head -n 5 "$file" | grep -q "Copyright"; then
            echo "  Copyright already exists"
            continue
        fi

        echo "  Adding copyright"

        {
            echo "$copyright"
            cat "$file"
        } > "$file.tmp"

        mv "$file.tmp" "$file"

    done < <(find "$rootDir" -type f)
done