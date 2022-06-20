# mock's User Guide

The simplest endpoint configuration we can define looks like this: 

```json
{
  "endpoints": [
    {
      "route": "foo/bar",
      "method": "POST",
      "response": {
        "foo": "bar"
      }
    }
  ]
}
```

A `POST` HTTP Request to `/foo/bar` will respond you with `{"foo":"bar"}`, as can be seen in the `response` endpoint configuration parameter above.

In the next sections we'll look at other ways how you can setup endpoints.

## Response with headers

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
+        "Some-Header-Key": "Some header value"
+      }
     }
   ]
}
```

## Response Status Code

By default, all responses' status code will be `200`. You can change it using the `response_status_code` option:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "POST",
       "response": {
         "foo": "bar"
       },
+      "response_status_code": 201 
     }
   ]
}
```


## File-based response content

In the earlier example, `response` is a JSON object containing the response JSON that you'll be responded with. However, as you setup complex APIs, your configuration file starts getting large and not easily readable. In the following example, we're setting the response content by referencing a file, thus leaving the configuration file more readable:

```json
{
  "endpoints": [
    {
      "route": "foo/bar",
      "method": "POST",
      "response": "file:./response_foobar.json"
    }
  ]
}
```

Given the configuration above, the `foo/bar` endpoint's response is defined in the `response_foobar.json` file.

## Conditional Response

You may want to define different responses for the same endpoint, based on certain conditions. The `response_if` parameter enables you to achieve that.

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "GET",
       "response": "Default response!",
+      "response_if": [
+        {
+          "response": "Hello world!",
+          "querystring_matches": [
+            {
+              "key": "foo",
+              "value": "bar"
+            }
+          ]
+        },
+        {
+          "response": "Hello galaxy!",
+          "querystring_matches": [
+            {
+              "key": "foo",
+              "value": "not_bar"
+            }
+          ]
         }
       ]
     }
   ]
}
```

In the configuration sample above, we have a single endpoint, `foo/bar`. There are two possible responses for this endpoint - if you call it with the `?foo=bar` querystring, the request response will be `Hello world!`. If however you use the `?foo=not_bar` querystring, the response will be `Hello galaxy!`.

> Note that, in the example above, even though we've added conditional responses, we still have a `response` like before - in the case where a request does not match any of the Response Conditions, the default `Default response!` response will be returned.
