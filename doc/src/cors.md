# Handling CORS

The `--cors` flag can be used when running *mock*. It will take care of setting up all the necessary headers in your API's Responses to enable browser clients to comunicate without problems:

```sh
$ mock serve --cors -c /path/to/your/config.json
```

