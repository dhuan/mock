# mock's User Guide

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**

- [Creating APIs](#creating-apis)
  - [Response with headers](#response-with-headers)
  - [Response Status Code](#response-status-code)
  - [File-based response content](#file-based-response-content)
  - [Conditional Response](#conditional-response)
    - [Condition Chaining](#condition-chaining)
    - [Headers in Conditional Responses](#headers-in-conditional-responses)
    - [Condition Options Reference](#condition-options-reference)
      - [`querystring_match`](#querystring_match)
      - [`form_match`](#form_match)
- [Test Assertions](#test-assertions)
  - [Which Request to assert against?](#which-request-to-assert-against)
  - [Assertion Chaining](#assertion-chaining)
  - [Assertion Options Reference](#assertion-options-reference)
    - [`form_match`](#form_match-1)
    - [`header_match`](#header_match)
    - [`json_body_match`](#json_body_match)
    - [`method_match`](#method_match)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

## Creating APIs

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

### Response with headers

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

### Response Status Code

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


### File-based response content

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

### Conditional Response

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

#### Condition Chaining

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

#### Headers in Conditional Responses

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

#### Condition Options Reference

In the earlier examples, we've seen that we can set Responses to be returned if a certain querystring matched, with the `querystring_match` Condition Option. There are, however, other Condition Options at your disposal for customizing your API.

You'll find all the available Condition Options in this section.

##### `querystring_match`

Matches against a Querystring in the Request. For example, a Request having the `?foo=bar` Querystring will be matched in the following condition:

```json
{
  "type": "querystring_match",
  "key": "foo",
  "value": "bar"
}
```

It's also possible to have multiple key/value pairs in the same condition. You will use the `key_values` field instead:

```json
{
  "type": "querystring_match",
  "key_values": {
    "some_key": "some value",
    "another_key": "another value"
  }
}
```

##### `form_match`

Matches against the Request's form-encoded data.

```json
{
  "type": "form_match",
  "key_values": {
    "some_key": "some value",
    "another_key": "another value"
  }
}
```

## Test Assertions

Besides enabling you to set-up APIs, mock provides you with methods to make assertions on how your endpoints were called.

Test Assertions are done by calling the `POST localhost:3000/__mock__/assert` endpoint.

In case you're new to the concept of automated tests and assertions - let's see what a very simple assertion looks like:

```json
{
  "route": "hello/world",
  "assert": {
    "type": "method_match",
    "value": "put"
  }
}
```

Or if we could say it in plain english: the endpoint `hello/world` was requested with the `put` method.

In case there was never a call to that particular endpoint, you would then get a response from mock indicating that no request has been made:

```json
{
  "validation_errors": [
    {
      "code": "no_call",
      "metadata": {}
    }
  ]
}
```

However in case a request had been made to that endpoint, with the, say, `POST` method, you would then get a different validation error, because you attempted to assert that it was called with the `PUT` method instead:

```json
{
  "validation_errors": [
    {
      "code": "method_mismatch",
      "metadata": {
        "method_expected": "put",
        "method_requested": "post"
      }
    }
  ]
}
```

> mock tells you whether the assertion passed or not by including "Validation Errors" into the `validation_errors` response field. Another indicative is the Response Status - `200` is success, `400` means your assertion failed.

With that we've seen a very simple assertion. There are other things that can be asserted in a HTTP Request, such as the header values passed, the body payload etc. [For a reference of all available assertion options, skip to this section.](#assertion-options-reference)

### Which Request to assert against?

By default, Assertions are based on the 1st Request. In cases where you want to assert against a Request other than the first, you'll use the `nth` Assertion Option.

```diff
 {
   "route": "foo/bar",
+  "nth": 2,
   "assert": {
     "type": "method_match",
     "value": "post"
   }
 }
```

### Assertion Chaining

Assertions can be combined with chaining options `and` and `or`. In the following assertion payload, we're extending the assertion we tried previously, asserting that our endpoint was called with the `{"foo":"bar"}` JSON payload.

```json
{
  "route": "hello/world",
  "assert": {
    "type": "method_match",
    "value": "post",
    "and": {
      "type": "json_body_match",
      "data": {
        "foo": "bar"
      }
    }
  }
}
```

In plain-english: assert that `hello/world` was requested with the `POST` method, **and** the `{"foo":"bar"}` JSON payload.

`or` can also be used for chaining assertions.

As shown in the example, chaining options are nested within a parent assertion. There's no limit to how many assertion chains you can make:

```json
{
  "route": "hello/world",
  "assert": {
    "type": "...",
    "value": "...",
    "and": {
      "type": "...",
      "value": "...",
      "and": {
        "type": "...",
        "value": "...",
        "or": {
          "type": "...",
          "value": "...",
          "and": {
            "type": "...",
            "value": "..."
          }
        }
      }
    }
  }
}
```

### Assertion Options Reference

#### `form_match`

Matches against the Request's Form Payload.

```json
{
  "type": "form_match",
  "key_values": {
    "some_form_key": "some_form_value",
    "another_form_key": "another_form_value"
  }
}
```

#### `header_match`

Matches against the Request's header.

```json
{
  "type": "header_match",
  "key_values": {
    "Some-header-key": "Some header value"
  }
}
```

#### `json_body_match`

The body payload que Request was called with.

```json
{
  "type": "json_body_match",
  "data": {
    "foo": "bar"
  }
}
```

#### `method_match`

The HTTP Method (Get, Post etc) the Request was called with.

```json
{
  "type": "method_match",
  "value": "post"
}
```
