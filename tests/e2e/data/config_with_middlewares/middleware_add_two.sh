printf "Two" >> $MOCK_RESPONSE_BODY

TMP=$(mktemp)
cat $MOCK_RESPONSE_STATUS_CODE > $TMP
NUM=$(cat $TMP)
printf $((NUM+1)) > $MOCK_RESPONSE_STATUS_CODE

echo 'New-Header-One: Value one' >> $MOCK_RESPONSE_HEADERS
