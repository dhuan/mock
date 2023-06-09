# Mock API Reference

Besides the custom endpoints defined in your configuration file, *mock* provides internal endpoints - these are identified by having a `__mock__` route prefix, such as the `/__mock__/assert` endpoint, which exists for making assertions. In this section you'll find out about each available internal endpoint.

## `POST __mock__/assert`

Makes Test Assertions, such as "endpoint X was called with Y payload.". The [Test Assertions Section](#test-assertions) dedicates to explaining all about assertions.

## `POST __mock__/reset`

Removes all Request Records that have been made so far. This has the same effect as stopping and starting *mock* over again. There are no parameters or payload fields to this endpoint.

