# mock

**mock** enables you to quickly setup a HTTP server for end-to-end tests. You can...

- Define endpoints and their respective responses;
- Make assertions on...
  - Whether a given endpoint was requested;
  - If a JSON payload body was passed correctly to a given endpoint;
  - If a header value was passed correctly;
  - And other useful things...

## Getting started

```sh
$ mock serve -c /path/to/your/config.json -p 4000
```

Below is a configuration file sample, defining two simple endpoints:

```
{
  "endpoints": [
    {
      "route": "foo/bar",
      "method": "POST",
      "response": {
        "foo": "bar"
      }
    },
    {
      "route": "hello/world",
      "method": "GET",
      "response": {
        "some_key": "some_value"
      }
    }
  ]
}
```

Let's now make a request to the `foo/bar` endpoint:

```sh
curl localhost:4000/foo/bar \
  -H 'Content-type: application/json' \
  -d '{"some_key":"some_value"}'
```

And then let's make assertions - Let's verify whether the `foo/bar` endpoint was called with the expected values:

```sh
curl -v http://localhost:4000/__mock__/assert -d @- <<EOF
{
    "route": "foo/bar",
    "body_json": {
      "some_key": "some_value"
    },
    "method": "put"
}
EOF
```

In the command above we're asserting that the `foo/bar` endpoint was called with the given payload, with the `put` method. Obviously, there's a problem with that assertion - the request we made previously was a `post` request, not `put`, therefore we get a response indicating so:

```sh
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

## Installing

mock is distributed as a single-file executable. Check the releases page and download the latest tarball.

## License

**mock** is licensed under MIT. For more information check the [LICENSE file.](LICENSE)
