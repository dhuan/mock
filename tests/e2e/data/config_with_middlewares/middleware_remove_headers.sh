TMP=$(mktemp)
cat $MOCK_RESPONSE_HEADERS | grep -iv one > $TMP
cat $TMP > $MOCK_RESPONSE_HEADERS
