echo Server Host: $MOCK_HOST >> $MOCK_RESPONSE_BODY
echo Request Host: $MOCK_REQUEST_HOST >> $MOCK_RESPONSE_BODY
echo URL: $MOCK_REQUEST_URL >> $MOCK_RESPONSE_BODY
echo Endpoint: $MOCK_REQUEST_ENDPOINT >> $MOCK_RESPONSE_BODY
echo Method: $MOCK_REQUEST_METHOD >> $MOCK_RESPONSE_BODY
echo Querystring: $MOCK_REQUEST_QUERYSTRING >> $MOCK_RESPONSE_BODY
echo Headers: >> $MOCK_RESPONSE_BODY
cat $MOCK_REQUEST_HEADERS >> $MOCK_RESPONSE_BODY
