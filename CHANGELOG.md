# Changelog

## Unreleased yet

IMPROVED:

- Avoid creating one temporary file for shell scripts for each request. An
  internal cache mechanism ensures that the amount of temporary files are
  reduced.

## 1.4.11

ADDED:

- get-payload: Add support for reading files from Multipart Form Requests.
  Previously only "fields" from Multipart could be read.

IMPROVED:

- get-route-param: Fail silently with Exit Code 1 instead of printing error
  message in case attempting to get Route Param for a route that has none.

FIXED:

- Requesting an unexisting endpoint while using `--cors` would panic.

## 1.4.10

FIXED:

- "Content-Length" header's value would be incorrect in case a Middleware
  changes a Base API's response body.

ADDED:

- Allow middlewares to modify CORS headers of Base APIs' Responses.

## 1.4.9

ADDED:

- Allow middlewares to execute on "not found" routes.

FIXED:

- Make `--route-match` flag work for Middlewares. This flag was documented in
  the guide but wasn't registered in the command-line parsing logic.

## 1.4.8

ADDED:

- Middlewares are now also executed for OPTIONS requests.
- Update OPTIONS HTTP Request logic so that the CORS headers are only included
  in the response when the "--cors" flag is used. If "--cors" is not used,
  then mock won't set OPTIONS routes for you automatically.

FIXED:

- Update `get-header` to return successful exit status code upon retrieving
  all headers.

## 1.4.7

ADDED:

- Add support for retrieving nested values from JSON payloads with `mock
  get-payload`. For example `mock get-payload users[0].name`.
- New ``--json`` option for ``mock write``.

## 1.4.6

ADDED:

- New commands: `mock get-header`, `mock set-status`, `mock get-payload`
- Add `--append` to `mock write`.

FIXED:

- BUG: Requests to file-server endpoints' index page would only be rendered if
  the URL contained a slash at the end. Not having a slash suffix would return
  404.

## 1.4.5

ADDED

- Directory navigation for file server.
- Add HTTP headers automatically for file-server routes, based on file
  extension
- New command: `mock get-query`
- Header `content-type: application/json` is now automatically added in case a
  response file is used with JSON extension.

## 1.4.4

ADDED

- New commands:
  - `mock set-header`.
  - `mock get-route-param`.
- `--regex` option for `mock replace`.

FIXED

- Middlewares could not use `mock write` and other helper commands.

IMPROVED

- Fail gracefully if more than 2 args are passed to `mock replace`.

## 1.4.3

ADDED

- Add ``mock write``.
- Add ``mock replace``.
- Support for regular expressions for `wipe-headers`, with `--regex
  <PATTERN>`.
- Decode gzipped responses retrieved from `forward`. Prior to this change,
  `$MOCK_RESPONSE_BODY` would contain gzip-compressed data if HTTP Response was
  encoded. Now you have access to the uncompressed data when manipulating the
  response body.

## 1.4.2

ADDED

- Add the `wipe-headers` command - a helper to be used within response shell
  scripts for removing undesired header.

FIXED

- Using `mock forward` with a Base API that uses HTTP/2 would fail due to a bug
  in the header parsing logic.

## 1.4.1

ADDED

- Automatically add `content-type: application/json` header if response is an
  array. Previously this automatic header addition would only occur if JSON
  body was an object at its root - array was not supported.

IMPROVED

- Prevent panic in case incorrectly formatted JSON is set as response.

FIXED

- `mock forward` would panic if Base API was defined without http protocol
  prefix.

## 1.4.0

IMPROVEMENTS

- Update CORS logic to prevent duplicated headers. If mock was running with a
  base api set which includes cors headers, then mock would respond with
  duplicated cors headers. Now mock overwrites any cors headers which were set
  by the Base API if executed with the `--cors` flag.

- Add `forward` command enabling Response Shell-Scripts to "forward" the
  current request to a Base API, and then modify the Base API's Response if
  desired. Check the "Base-API" section in the manual for more details.

## 1.3.1

IMPROVEMENTS:

- If started without specifying an http port, then a random available port will
  be used, instead of trying port 3000.
- Add validation for "Base API" value. If you try to set it as an invalid
  hostname or domain, *mock* will fail gracefully.

## 1.3.0

CHANGED

- Allow mock to be started only with a Base API - no endpoints set, acting only
  as a proxy.

FIXED

- Base APIs could only be used through command-line flag (`--base`), although
  it was documented that the `base` json config parameter was supported as
  well. This has been fixed - the feature is now usable through cmd flag and
  json config.

## 1.2.0

ADDED

- Support for Base APIs
- New environment variable for capturing individual querystring parameters:
  `$MOCK_REQUEST_QUERYSTRING_FOOBAR`
- Enable plaintext responses to use *mock* variables, such as:

```
mock serve --route foo/bar --response 'Url: ${MOCK_REQUEST_URL}'
```
- Enable individual request headers to be read through environment variables
  such has `MOCK_REQUEST_HEADER_FOO_BAR`


## 1.1.0

ADDED

- Add support for "conditions" for Middlewares.
- Enable Middlewares' `exec` commands to include any shell operators like
  pipes, output redirection, etc.
- New condition options added: `querystring_match_regex`, `querystring_exact_match_regex`.

## 1.0.0

BREAKING CHANGES

- Request Handlers (shell scripts and executables) now need to write to the
  `$MOCK_RESPONSE_BODY` environment variable in order to write an HTTP
  Response's body instead of outputting to stdout.
- Interfaces for assertions have changed for better readability: `assert` field
  has been renamed to `condition`. Check the documentation on the assertion
  sections before upgrading.

ADDED

- Middlewares support;
- New environment variable for Request/Middleware Handlers: `MOCK_REQUEST_NTH`;

## 0.8.1

ADDED

- New condition option: `route_param_match`

## 0.8.0

ADDED

- `nth` condition option: Conditional responses can be set based on their position in the request history.

## 0.7.0

ADDED

- Responses with `exec:<SHELL COMMAND>` is now supported;
- Endpoints can now be defined through command-line parameters, such as: `--route foo/bar --method post --status-code 201 --response "Hello world!"`;

CHANGED

- The `--config` command-line parameter is no longer mandatory since now endpoints can be defined without configuration files;

FIXED

- Defining responses referenced through files with absolute path fails. (with relative file paths no issues, only absolute);

## 0.6.0

ADDED

- Responses can now read environment variables. Previously only shell-script responses had that ability - now any kind of response, either file or static text can achieve the same. Check *Reading Environment Variables* in the User Guide.
- New variable added to read current request's host `MOCK_REQUEST_HOST`.

## 0.5.0

Breaking changes

- "Endpoint Parameters" has been renamed to "Route Parameters"

Example - reading a parameter named `foo`:

Before: `MOCK_REQUEST_ENDPOINT_PARAM_FOO`

Now: `MOCK_ROUTE_PARAM_FOO`

Check the User Guide for more details.

Features & enhancements

- Route Parameters can be captured in the Response string. Before, the parameters could only be read by Shell Scripts Responses. A response can now be set as follows:

```json
{
  "endpoints": [
    {
      "route": "book/{book_name}",
      "method": "GET",
      "response": "file:./books/${book_name}.txt"
    }
  ]
}
```

- Static files support, with `fs:./path/to/files`;
- Endpoints can be configured without any HTTP Method - it will default to `GET`;

## 0.4.0

Additions:

- New Request Handler Variable added: `MOCK_HOST` for retrieving the current host that the Mock server is listening on.
- Enable JSON Responses to include environment variables, which previously could only be read by Shell Script responses.

Bugs fixed:

- Trying to assert with "Json Body" on a Request that didn't have any payload would result in 500 Status Code API Error;

## 0.3.0

- Endpoint Routes can now contain wildcards or placeholder variables;
- Shell Script Response Handlers now support params, such as: `"response": "sh:./my_shell_script.sh some_param another_param"`;
- `--delay` option added enabling you to simulate slow APIs;

## 0.2.0

- Endpoint responses from shell scripts are now supported (with `sh:some_handler.sh`). Read the User Guide for more details;
- `--cors` option added to facilitate usage with webapps;

Minor stuff:

- When trying to reference response files that do not exist, *mock* now prevents starting and shows error, failing gracefully.

## 0.1.4

- The *mock* installable Go library has existed before this release but now it is documented in the User Guide. Add `github.com/dhuan/mock/pkg/mock` to your Go project and write tests with it.
- Helper function `ToReadableError()` added in the library to stringify a group of Validation Errors.
- Bug fixed - HTTP Method value in Assertion now works independently of case sensitiveness.

Minor stuff:

- Fail gracefully if given configuration file does not exist or/and not readable.

## 0.1.3

This release fixes a bug in the `json_body_match` condition option.

## 0.1.2

This release doesn't have significant changes.

## 0.1.1

Features:

- `querystring_exact_match` Assertion Matcher added;

## 0.1.0

Features:

- The `querystring` assertion matcher was added. You can now assert that a request was made with the desired Querystring values and keys.

General improvements and stability:

- Log messages are shown with timestamps.
- Proper error handling when unable to start up the server.

Bugs fixed:

- Trying to make assertions with a `nth` out of range would result in panicking the server. Mock now returns a proper validation error on the assert request indicating that the given `nth` is out of range.
