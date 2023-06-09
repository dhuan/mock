# Assertion Chaining

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

