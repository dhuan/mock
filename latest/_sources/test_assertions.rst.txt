Test Assertions
===============

Besides enabling you to set-up APIs, mock provides you with methods to
make assertions on how your endpoints were called.

Test Assertions are done by calling the
``POST localhost:3000/__mock__/assert`` endpoint.

In case you’re new to the concept of automated tests and assertions -
let’s see what a very simple assertion looks like:

.. code:: json

   {
     "route": "hello/world",
     "condition": {
       "type": "method_match",
       "value": "put"
     }
   }

Or if we could say it in plain english: the endpoint ``hello/world`` was
requested with the ``put`` method.

In case there was never a call to that particular endpoint, you would
then get a response from mock indicating that no request has been made:

.. code:: json

   {
     "validation_errors": [
       {
         "code": "no_call",
         "metadata": {}
       }
     ]
   }

However in case a request had been made to that endpoint, with the, say,
``POST`` method, you would then get a different validation error,
because you attempted to assert that it was called with the ``PUT``
method instead:

.. code:: json

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

..

   mock tells you whether the assertion passed or not by including
   “Validation Errors” into the ``validation_errors`` response field.
   Another indicative is the Response Status - ``200`` is success,
   ``400`` means your assertion failed.

With that we’ve seen a very simple assertion. There are other things
that can be asserted in a HTTP Request, such as the header values
passed, the body payload etc. `For a reference of all available
condition options, skip to this section. <conditions_reference.html>`__

Which Request to assert against?
--------------------------------

By default, Assertions are based on the 1st Request. In cases where you
want to assert against a Request other than the first, you’ll use the
``nth`` Assertion Option.

.. code:: diff

    {
      "route": "foo/bar",
   +  "nth": 2,
      "condition": {
        "type": "method_match",
        "value": "post"
      }
    }

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   assertion_chaining
   mock_package
