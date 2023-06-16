TMP=$(mktemp -d)
echo $TMP

php -S localhost:8080 -t "$TMP" &

find doc/source -type f | entr -s 'make -s doc_build | tar xz -C '"$TMP"
