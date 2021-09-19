#!/usr/bin/env sh

# Reformat all feature files using https://github.com/antham/ghokin

exitCode=0
for file in "$@"; do
    output="$(ghokin check "$file" 2>&1)"
    outputCode="$?"
    if [ ${outputCode} -ne 0 ]; then
        echo "$output"
        ghokin fmt replace "$file"
        exitCode=1
    fi
done

exit ${exitCode}
