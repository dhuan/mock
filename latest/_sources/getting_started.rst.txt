.. _getting_started:

Getting started
===============

Let's look into the basics.

The basics - "hello world" with mock
------------------------------------

The simplest way of defining an endpoint's response is simply using the
``response`` field within an endpoint object:

.. code:: js

    {
      "endpoints": [
        {
          "route": "foo/bar",
          "response": "Hello, world!"
        }
      ]
   }

The exact same can also be accomplished with just command-line parameters:

.. code:: sh

    $ mock serve \
      --route "foo/bar" \
      --response 'Hello, world!'

For simple APIs, static response such as the ones above might be enough. But we're not limited to that.

`mock` sets no boundaries when it comes to building APIs. You have all the
capability of the shell environment at your disposal for building APIs.

mock as an API framework
------------------------

`mock` lets you use shell scripts as "response handlers" for your endpoints.

.. code:: sh

    $ mock serve \
      --route "/users" \
      --method "GET" \
      --shell-script 'get_users.sh' \
      --route "/users" \
      --method "POST" \
      --shell-script 'create_user.sh'

Above we have an API with a ``/users`` route, that can be requested with either
``GET`` or ``POST`` methods. The shell script files ``get_users.sh`` and
``create_user.sh`` will be executed respectively, acting as response handlers
for these endpoints.

Inside your response handlers you can use `mock` to build HTTP Responses.

``mock write`` is one such utility. It sends data to the requesting client.
Let's see a simple example:

.. code:: sh

   # get_users.sh

   USERS=$(mysql ... -c "SELECT * FROM users")

   # convert the user sql result to JSON somehow...

   echo "${USERS}" | mock write

You can read more about shell scripts and shell utilities in the following
links:

- :ref:`Response from shell scripts<shell_script_responses>`
- :ref:`Shell utilities<shell_utils>`

HTTP Headers
------------

The optional ``response_headers`` endpoint parameter will add headers to
a endpoint’s response:

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
          "response": {
            "foo": "bar"
          },
   +      "response_headers": {
   +        "Some-Header-Key": "Some header value",
   +        "Another-Header-Key": "Another header value"
   +      }
        }
      ]
   }

To add response headers to an endpoint using command-line parameters:

.. code:: diff

    $ mock serve \
      --route "foo/bar" \
      --response '{"foo":"bar"}' \
   +  --header "Some-Header-Key: Some header value" \
   +  --header "Another-Header-Key: Another header value"

If you're using :ref:`shell scripts as response handlers
<shell_script_responses>`, then setting up headers is quite easy.
Let's use :ref:`mock set-header <shell_utils_set_header>` to add a ``Foo: bar`` HTTP Header to our API response:

.. code:: diff

    $ mock serve \
      --route "foo/bar" \
   +  --exec 'mock set-header foo bar'

HTTP Status Code
----------------

By default, all responses’ status code will be ``200``. You can change
it using the ``response_status_code`` option:

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
          "method": "POST",
          "response": {
            "foo": "bar"
          },
   +      "response_status_code": 201 
        }
      ]
   }

To add response status codes to an endpoint using command-line
parameters:

.. code:: diff

    $ mock serve \
      --route "foo/bar" \
      --method "POST" \
      --response '{"foo":"bar"}' \
   +  --status-code 201
