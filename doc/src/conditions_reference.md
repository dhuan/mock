# Conditions Reference

When making test assertions with *mock*, *Conditions* enable you to express how you expect a given Request to have been made. *Conditions* are also used when defining [*Conditional Responses*.](conditional_responses.md)

In this section you will find a reference of all available Conditions.

## `querystring_match`

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

## `querystring_exact_match`

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

## `json_body_match`

Matches against the JSON body payload que Request was called with.

```json
{
  "type": "json_body_match",
  "key_values": {
    "foo": "bar"
  }
}
```

## `form_match`

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

## `header_match`

Matches against the Request's header.

```json
{
  "type": "header_match",
  "key_values": {
    "Some-header-key": "Some header value"
  }
}
```

## `method_match`

Matches against the HTTP Method (Get, Post etc) the Request was called with.

```json
{
  "type": "method_match",
  "value": "post"
}
```

## `route_param_match`

Matches against the Route Param in the requested endpoint.

```json
{
  "type": "route_param_match",
  "key": "some_param_name",
  "value": "some_value"
}
```

## `nth`

Matches if the current request is nth on the request history. Note that both route and method must match. In the example below, a match will occur only if the request is the 2nd made so far to the server.

```json
{
  "type": "nth",
  "value": 2
}
```

It's also possible to match all subsequent requests after a given number, just add a "+" (plus) sign after the number (note also that to accomplish this, the value must be defined as a string). For example, let's match all requests starting from the second onwards:

```json
{
  "type": "nth",
  "value": "2+"
}
```

