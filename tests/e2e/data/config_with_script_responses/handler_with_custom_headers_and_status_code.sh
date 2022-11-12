echo Hello world!

cat <<EOF > $MOCK_RESPONSE_HEADERS
Some-Header-Key: Some Header Value
Another-Header-Key: Another Header Value
EOF

echo 201 > $MOCK_RESPONSE_STATUS_CODE
