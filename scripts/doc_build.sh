TMP=$(mktemp -d)

cp -r doc "$TMP"

cp CHANGELOG.md "$TMP""/doc/src/changelog.md"

cd "$TMP""/doc"

mdbook build

cd book

tar cz ./*
