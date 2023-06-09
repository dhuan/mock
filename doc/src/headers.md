# Response with headers

The optional `response_headers` endpoint parameter will add headers to a endpoint's response:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "POST",
       "response": {
         "foo": "bar"
       },
+      "response_headers": {
+        "Some-Header-Key": "Some header value",
+        "Another-Header-Key": "Another header value"
+      }
     }
   ]
}
```

To add response headers to an endpoint using command-line parameters:

```diff
 $ mock serve \
   --route "foo/bar" \
   --method "POST" \
   --response '{"foo":"bar"}' \
+  --header "Some-Header-Key: Some header value" \
+  --header "Another-Header-Key: Another header value"
```

