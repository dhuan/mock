# Changelog

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
