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
  - [Handling CORS](#handling-cors)
- [Test Assertions](#test-assertions)
  - [Which Request to assert against?](#which-request-to-assert-against)
  - [Assertion Chaining](#assertion-chaining)
- [Test Assertions with *mock*'s Go package](#test-assertions-with-mocks-go-package)
- [Conditions Reference](#conditions-reference)
  - [`querystring_match`](#querystring_match)
  - [`querystring_exact_match`](#querystring_exact_match)
  - [`json_body_match`](#json_body_match)
  - [`form_match`](#form_match)
  - [`header_match`](#header_match)
  - [`method_match`](#method_match)
- [Mock API Reference](#mock-api-reference)
  - [`POST __mock__/assert`](#post-__mock__assert)
  - [`POST __mock__/reset`](#post-__mock__reset)

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

In the examples above, we've seen that we can set Responses to be returned if a certain querystring matched, with the `querystring_match` Condition Option. There are, however, other Condition Options at your disposal for customizing your API. [Read the Condition Reference for a list of all available Conditions.](#conditions-reference)

### Handling CORS

The `--cors` flag can be used when running *mock*. It will take care of setting up all the necessary headers in your API's Responses to enable browser clients to comunicate without problems:

```sh
$ mock serve --cors -c /path/to/your/config.json
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

With that we've seen a very simple assertion. There are other things that can be asserted in a HTTP Request, such as the header values passed, the body payload etc. [For a reference of all available condition options, skip to this section.](#conditions-reference)

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
      "key_values": {
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

## Test Assertions with *mock*'s Go package

In the previous section we've seen how to make test assertions by means of HTTP requests. With that we've seen how *mock* is designed to be language-agnostic - no matter what programming language you're using for your E2E tests, *mock* can easily be integrated because HTTP requests are all that's needed for making test assertions. But we're not limited to HTTP requests only, when making assertions. In this section we'll learn how to use *mock*'s Go package, which enables you to achieve the same but without requiring to write requests by hand.

Let's take as an example, a test assertion in its plain request format, asserting that a request was made to `foo/bar` with the `POST` method.

```sh
curl -v -X POST "localhost:4000/__mock__/assert" -d @- <<EOF
{
  "route": "foo/bar",
  "assert": {
    "type": "method_match",
    "value": "post"
  }
}
EOF
```

Let's now convert it to Go code - what does a test case (using *mock*) look like in Go?

```go
package my_test

import (
	"github.com/dhuan/mock/pkg/mock"
	"testing"
)

func Test_FooBarShouldBeRequested(t *testing.T) {
	mockConfig := &mock.MockConfig{Url: "localhost:4000"}

	validationErrors, err := mock.Assert(mockConfig, &mock.AssertConfig{
		Route: "foo/bar",
		Assert: &mock.Condition{
			Type:  mock.ConditionType_MethodMatch,
			Value: "post",
		},
	})

	if err != nil {
		t.Error(err)
	}

	if len(validationErrors) > 0 {
		t.Error(mock.ToReadableError(validationErrors))
	}
}
```

Just like you get a response containing Validation Errors when using the HTTP-request approach, in Go the Validation Errors are returned from the `mock.Assert(...)` call.

> Note that with *mock*'s Go package, we're simply executing assertions. The actual *mock server instance* is supposed to be running and started before your test script starts.

A few things to be noted regarding the Go code snippet above:

- Prior to making assertions, you need to tell the *mock* library what network host+port *mock* is running at, which is done with `&mock.MockConfig{Url: "localhost:4000"}`
- Besides the `validationErrors` returned from `mock.Assert(...)`, we still get a 2nd return value of type `error`. This error is not related to *mock*'s Validation Errors. This error can be something like if HTTP failure in case *mock* is not running on the network port you set it to. It's important to check and fail the test if `err` is not `nil` (as shown in the example), otherwise it will seem as if your test passed because there are no Validation Errors but an actual error occurred.
- If `validationErrors` is an empty slice and `err` is nil, then your assertion passed successfully.

With that we covered basic assertions. Let's see now a more complex kind of assertion, using *Assertion Chaining*:

```diff
 validationErrors, err := mock.Assert(mockConfig, &mock.AssertConfig{
     Route: "foo/bar",
     Assert: &mock.Condition{
         Type:  mock.ConditionType_MethodMatch,
         Value: "post",
+        And: &mock.Condition{
+            Type: mock.ConditionType_JsonBodyMatch,
+            KeyValues: map[string]interface{}{
+                "some_key": "some_value",
+            },
+        },
     },
 })
```

## Conditions Reference

When making test assertions with *mock*, *Conditions* enable you to express how you expect a given Request to have been made. *Conditions* are also used when defining [*Conditional Responses*.](#conditional-response)

In this section you will find a reference of all available Conditions.

### `querystring_match`

Matches against a Querystring in the Request. For example, a Request having the `?foo=bar` Querystring will be matched in the following condition:

```json
{
  "type": "querystring_match",
  "key": "foo",
  "value": "bar"
}
```

You can also use `key_values` and set multiple pairs:

```json
{
  "type": "querystring_match",
  "key_values": {
    "some_key": "some value",
    "another_key": "another value"
  }
}
```

### `querystring_exact_match`

Matches against Querystring values, like `querystring_match`. The difference being that it matches only if the Request's Querystring contains only the specified Querystrings and no other.

```json
{
  "type": "querystring_exact_match",
  "key_values": {
    "some_key": "some value",
    "another_key": "another value"
  }
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

### `json_body_match`

Matches against the JSON body payload que Request was called with.

```json
{
  "type": "json_body_match",
  "key_values": {
    "foo": "bar"
  }
}
```

### `form_match`

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

### `header_match`

Matches against the Request's header.

```json
{
  "type": "header_match",
  "key_values": {
    "Some-header-key": "Some header value"
  }
}
```

### `method_match`

Matches against the HTTP Method (Get, Post etc) the Request was called with.

```json
{
  "type": "method_match",
  "value": "post"
}
```


## Mock API Reference

Besides the custom endpoints defined in your configuration file, *mock* provides internal endpoints - these are identified by having a `__mock__` route prefix, such as the `/__mock__/assert` endpoint, which exists for making assertions. In this section you'll find out about each available internal endpoint.

### `POST __mock__/assert`

Makes Test Assertions, such as "endpoint X was called with Y payload.". The [Test Assertions Section](#test-assertions) dedicates to explaining all about assertions.

### `POST __mock__/reset`

Removes all Request Records that have been made so far. This has the same effect as stopping and starting *mock* over again. There are no parameters or payload fields to this endpoint.
