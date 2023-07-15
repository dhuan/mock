Test Assertions with *mock*\ ’s Go package
==========================================

In the previous section we’ve seen how to make test assertions by means
of HTTP requests. With that we’ve seen how *mock* is designed to be
language-agnostic - no matter what programming language you’re using for
your E2E tests, *mock* can easily be integrated because HTTP requests
are all that’s needed for making test assertions. But we’re not limited
to HTTP requests only, when making assertions. In this section we’ll
learn how to use *mock*\ ’s Go package, which enables you to achieve the
same but without requiring to write requests by hand.

Let’s take as an example, a test assertion in its plain request format,
asserting that a request was made to ``foo/bar`` with the ``POST``
method.

.. code:: sh

   curl -v -X POST "localhost:4000/__mock__/assert" -d @- <<EOF
   {
     "route": "foo/bar",
     "assert": {
       "type": "method_match",
       "value": "post"
     }
   }
   EOF

Let’s now convert it to Go code - what does a test case (using *mock*)
look like in Go?

.. code:: go

   package my_test

   import (
       "github.com/dhuan/mock/pkg/mock"
       "testing"
   )

   func Test_FooBarShouldBeRequested(t *testing.T) {
       mockConfig := &mock.MockConfig{Url: "localhost:4000"}

       validationErrors, err := mock.Assert(mockConfig, &mock.AssertOptions{
           Route: "foo/bar",
           Condition: &mock.Condition{
               Type:  mock.ConditionType_MethodMatch,
               Value: "post",
           },
       })

       if err != nil {
           t.Error(err)
       }

       if len(validationErrors) > 0 {
           t.Error(mock.ToReadableError(validationErrors))
       }
   }

Just like you get a response containing Validation Errors when using the
HTTP-request approach, in Go the Validation Errors are returned from the
``mock.Assert(...)`` call.

.. note::

   Note that with *mock*\ ’s Go package, we’re simply executing
   assertions. The actual *mock server instance* is supposed to be
   running and started before your test script starts.

A few things to be noted regarding the Go code snippet above:

-  Prior to making assertions, you need to tell the *mock* library what
   network host+port *mock* is running at, which is done with
   ``&mock.MockConfig{Url: "localhost:4000"}``
-  Besides the ``validationErrors`` returned from ``mock.Assert(...)``,
   we still get a 2nd return value of type ``error``. This error is not
   related to *mock*\ ’s Validation Errors. This error can be something
   like if HTTP failure in case *mock* is not running on the network
   port you set it to. It’s important to check and fail the test if
   ``err`` is not ``nil`` (as shown in the example), otherwise it will
   seem as if your test passed because there are no Validation Errors
   but an actual error occurred.
-  If ``validationErrors`` is an empty slice and ``err`` is nil, then
   your assertion passed successfully.

With that we covered basic assertions. Let’s see now a more complex kind
of assertion, using *Assertion Chaining*:

.. code:: diff

    validationErrors, err := mock.Assert(mockConfig, &mock.AssertOptions{
        Route: "foo/bar",
        Condition: &mock.Condition{
            Type:  mock.ConditionType_MethodMatch,
            Value: "post",
   +        And: &mock.Condition{
   +            Type: mock.ConditionType_JsonBodyMatch,
   +            KeyValues: map[string]interface{}{
   +                "some_key": "some_value",
   +            },
   +        },
        },
    })
