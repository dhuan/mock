# Changelog

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
