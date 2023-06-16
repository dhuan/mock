TMP=$(mktemp -d)
TMP_OUTPUT=$(mktemp)

cp -r doc "$TMP"
cp 'CHANGELOG.md' "$TMP"'/.'

cd "$TMP""/doc"

pandoc --from=markdown --to=rst --output='source/changelog.rst' ../CHANGELOG.md

if ! make clean > "$TMP_OUTPUT";
then
    echo $TMP_OUTPUT
fi

if ! make html > "$TMP_OUTPUT";
then
    echo $TMP_OUTPUT
fi

cd build/html

tar cz ./*
