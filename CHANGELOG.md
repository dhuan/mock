# Changelog

## master (not released yet)

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
