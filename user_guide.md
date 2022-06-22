# mock's User Guide

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [mock's User Guide](#mocks-user-guide)
  - [Response with headers](#response-with-headers)
  - [Response Status Code](#response-status-code)
  - [File-based response content](#file-based-response-content)
  - [Conditional Response](#conditional-response)
    - [Condition Chaining](#condition-chaining)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


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
+          "conditions": [
+            {
+              "type": "querystring_match",
+              "key": "foo",
+              "value": "bar"
+            }
+          ]
+        }
       ]
     }
   ]
}
```

In the configuration sample above, we have a single endpoint, `foo/bar`. There are two possible responses for this endpoint - if you call it with the `?foo=bar` querystring, the request response will be `Hello world!`. If however you use the `?foo=not_bar` querystring, the response will be `Hello galaxy!`.

> Note that, in the example above, even though we've added conditional responses, we still have a `response` like before - in the case where a request does not match any of the Response Conditions, the default `Default response!` response will be returned.

### Condition Chaining

In the previous example we defined a Response with a very simple querystring-based condition. Next we'll look at how to define more complex conditions,  with condition chaining, which is possible with the `and` and `or` combination options. We'll extend the previous condition example: besides the `foo=bar` querystring value, it will also be necessary that the `hello=world` querystring is present in the request.

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "GET",
       "response": "Default response!",
       "response_if": [
         {
           "response": "Hello world!",
           "conditions": [
             {
               "type": "querystring_match",
               "key": "foo",
               "value": "bar",
+              "and": {
+                "type": "querystring_match",
+                "key": "hello",
+                "value": "world"
+              }
             }
           ]
         }
       ]
     }
   ]
}
```

Now, the `Hello world!` Response will only be returned if the request has the following querystring values: `foo=bar&hello=world`.

Besides the `and` option, you can also use `or`.

There's no limit to how deep you can nest a chain of conditions.
