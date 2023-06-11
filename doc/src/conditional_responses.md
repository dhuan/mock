# Conditional Response

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
+          "condition": {
+            "type": "querystring_match",
+            "key": "foo",
+            "value": "bar"
+          }
+        }
+      ]
     }
   ]
 }
```

In the configuration sample above, we have a single endpoint, `foo/bar`. There are two possible responses for this endpoint - if you call it with the `?foo=bar` querystring, the request response will be `Hello world!`. If however you use the `?foo=not_bar` querystring, the response will be `Hello galaxy!`.

> Note that, in the example above, even though we've added conditional responses, we still have a `response` like before - in the case where a request does not match any of the Response Conditions, the default `Default response!` response will be returned.

## Condition Chaining

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
           "condition": {
             "type": "querystring_match",
             "key": "foo",
             "value": "bar",
+            "and": {
+              "type": "querystring_match",
+              "key": "hello",
+              "value": "world"
+            }
           }
         }
       ]
     }
   ]
 }
```

Now, the `Hello world!` Response will only be returned if the request has the following querystring values: `foo=bar&hello=world`.

Besides the `and` option, you can also use `or`.

There's no limit to how deep you can nest a chain of conditions.

## Headers in Conditional Responses

Conditional Response Objects *can* have Headers set as well. Use the `response_headers` option:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "GET",
       "response": "Default response!",
       "response_headers": {
         "Header-Foo": "Foobar!"
       },
       "response_if": [
         {
           "response": "Hello world!",
+          "response_headers": {
+            "Some-Header-Key": "Some header value"
+          },
           "condition": {
             /* ... */
           }
         }
       ]
     }
   ]
 }
```

If you don't set a `response_header` to a Conditional Response, the response will not have any headers, even if a `response_headers` field exists in the main Response.

If you'd like Conditional Responses to inherit the Headers from the main Response, use the `response_headers_base`:

```diff
 {
   "endpoints": [
     {
       "route": "foo/bar",
       "method": "GET",
       "response": "Default response!",
       "response_headers": {
         "Header-Foo": "Foobar!"
       },
+      "response_headers_base": {
+        "Some-base-header": "Some value for the base header"
+      },
       "response_if": [
         {
           "response": "Hello world!",
           "response_headers": {
             "Some-Header-Key": "Some header value"
           },
           "condition": {
             /* ... */
           }
         }
       ]
     }
   ]
 }
```

With the example configuration above, a Request resolving to the Conditional Response would result in having the following headers:

```
Some-base-header: Some value for the base header
Some-Header-Key: Some header value
```

In the examples above, we've seen that we can set Responses to be returned if a certain querystring matched, with the `querystring_match` Condition Option. There are, however, other Condition Options at your disposal for customizing your API. [Read the Condition Reference for a list of all available Conditions.](conditions_reference.md)

