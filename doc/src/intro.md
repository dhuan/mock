# Welcome to mock's documentation!

*mock* enables you to quickly set up HTTP servers for end-to-end tests.

- Define endpoints and their respective responses through easy syntax;
- Make assertions on...
  - Whether a given endpoint was requested;
  - If a JSON payload body was passed correctly to a given endpoint;
  - If a header value was passed correctly;
  - And other useful things...

```sh
$ mock serve --port 3000 \
  --route 'time_in/{country}' \
  --method GET \
  --exec 'zdump ${country}' \
  --route 'whois/{domain}' \
  --method GET \
  --exec 'whois ${domain}'
```

Run the example command the above and try these URLs in your browser or any preferred HTTP client: `http://localhost:3000/time_in/Japan` and `http://localhost:3000/whois/google.com`

## Quick links

- [Download *mock* for Linux](__DOWNLOAD_LINK_LINUX__)
- [Download *mock* for MacOS](__DOWNLOAD_LINK_MACOS__)
- [Releases](https://github.com/dhuan/mock/releases)
- [*mock*'s source code](https://github.com/dhuan/mock)
- [Report bugs](https://github.com/dhuan/mock/issues)

## Read further...

The core functionalities of *mock* are documented each in their respective sections. Read further to learn:
- [Creating APIs](apis.md)
- [Test Assertions](test_assertions.md)

## License

*mock* is licensed under MIT. For more information check the [LICENSE file.](LICENSE)
