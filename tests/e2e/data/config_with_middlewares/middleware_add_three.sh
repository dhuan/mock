printf "Three" >> $MOCK_RESPONSE_BODY

TMP=$(mktemp)
cat $MOCK_RESPONSE_STATUS_CODE > $TMP
NUM=$(cat $TMP)
printf $((NUM+2)) > $MOCK_RESPONSE_STATUS_CODE

echo 'New-Header-Two: Value two' >> $MOCK_RESPONSE_HEADERS
