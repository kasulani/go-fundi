#!/usr/bin/env sh

# Audit git repos for secrets using https://github.com/zricethezav/gitleaks

gitleaks --verbose --pretty --redact --depth 100 --repo "/go-fundi"

if [ $? -eq 1 ]; then
    cat <<\EOF
Alert: Possible exposure of sensitive information was identified in your changes.
You must review your changes and make sure that no passwords, keys, tokens
and/or secrets are commited.
If you think that this is a false positive, please report it
to security@hellofresh.com with the link to your commit.
We appreciate your help.
HelloFresh Information Security Team
EOF
    exit 1
fi
