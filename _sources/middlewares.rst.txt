.. _middlewares:

Middlewares
===========

Middlewares enable you to execute any kind of program logic before responses
are sent out to clients. Any program that's executable from shells can be a
middleware. All properties of a response can be manipulated by middlewares,
such as its body, headers, status code, etc.

Setting up Middlewares
~~~~~~~~~~~~~~~~~~~~~~

Following is a *mock* configuration file containing middlewares:

.. code:: json

   {
    "middlewares": [
        {
            "exec": "sh /path/to/middleware/script.sh",
            "route_match": "*"
        }
    ],
    "endpoints": [
        // ...
    ]
   }

Let's go over each the fields shown in the ``middleware`` objects:

- ``exec``: A shell command that will be executed which can perform a
  middleware operation, such as changing request/response. Anything that's
  valid in the shell will be valid here, such as ``sh path/to/some/script.sh``
  or ``/path/to/some/program``. Shell operators can also be used here such as
  pipes and output redirection.
- ``route_match``: Optional. If not set or set as ``*`` then all Requests will
  be processed for the given Middleware. For filtering only desired requests,
  set the value as a Regular Expression, for example ``foo/[a-z]{1,}`` which will
  match against requests such as ``foo/bar``.

Middlewares can also be set from command-line parameters as an alternative to
configuration files:

.. code:: diff

   $ mock serve \
     --route foo/bar
     --response "Hello world!"
  +  --middleware "sh path/to/some/script.sh"
  +  --route-match 'foo/bar'
  +  --middleware "sh path/to/another/script.sh"

.. note::

   In the example above, note that the 2nd middleware does not use the
   ``--route-match`` option, which will result in that middleware being
   executed for all requests. The 1st middleware in the example however uses
   the route matching option.

Examples of Middlewares
~~~~~~~~~~~~~~~~~~~~~~~

.. note::

   The examples below are provided as shell scripts programs. But remember that
   middlewares are not limited to being shell scripts only. Any executable
   program at all can be used as a middleware.

Modify Response Body
--------------------

To kick-off the middleware examples, let's set up a very simple middleware that
will uppercase all response text.

.. code:: sh

   $ mock serve -p 3000 \
       --middleware 'awk '"'"'{print toupper($0)}'"'"' $MOCK_RESPONSE_BODY | mock write' \
       --route foo/bar \
       --response 'Hello, world!'

The middleware above simply pipes all response data to ``awk {print
toupper($0)}``, uppercasing all text.

Let's now request this `mock API` and find out the results:

.. code:: sh

   $ curl localhost:3000/foo/bar

   // Prints out: HELLO, WORLD!

.. note::

   To modify the response in the example above, we've used ``mock write``. If
   that's new to you :ref:`read more about it here. <shell_utils_write>`

Adding new headers before sending response to client
----------------------------------------------------

The following middleware adds a header to all endpoints:

.. code:: sh

    $ mock serve -p 3000 \
        --middleware 'mock set-header some-header some-value' \
        --route foo/bar \
        --response 'Hello, world!'

Let's request our `mock API` and find out if the header was used:

.. code:: sh

    $ curl -v localhost:3000/foo/bar

    // Prints out:

    > GET /foo/bar HTTP/1.1                                                                                                                                       
    > Host: localhost:3000                                                                                                                                        
    >                                                                                                                                                             
    < HTTP/1.1 200 OK                                                                                                                                             
    < Some-Header: some-value                                                                                                                                     
    <                                                                                                                                                             
    { [13 bytes data]                                                                                                                                             
    Hello, world!                             

As we can see the ``some-header`` header was included in the response, thanks
to the middleware. Note the usage of ``-v`` in CURL otherwise we could not have
seen the response headers.

.. note ::

    In the example above we used ``mock set-header``. :ref:`Read more about it
    here. <shell_utils_set_header>`

Environment Variables for Middlewares
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

Middleware Handlers are provided with a set of environment variables with
useful information about the request being processed, and also files that can
be written to to customize your API behavior.

The following variables hold file paths that can be written to in order to
customise responses:

- ``MOCK_RESPONSE_BODY``: A file that can be written to in order to modify the
  HTTP Response before handing it to the client. This file already contains the
  response body defined by your API configuration for the given endpoint.
- ``MOCK_RESPONSE_HEADERS``: A file that can be written to in order to modify
  the HTTP Headers. The headers defined in your configuration's endpoint are
  included in this file at the moment the middleware is executed.
- ``MOCK_RESPONSE_STATUS_CODE``: A file that can be written to in order to
  modify the HTTP Status code.

Route Parameters can also be read. For example if an endpoint exists with its
route set as ``foo/bar/{some_param}``, middlewares can read them through
environment variables such as ``MOCK_ROUTE_PARAM_SOME_PARAM``

As Middlewares are executed for both requests with valid routes and requests
for which a route hasn't been set, you may want to know inside your handler
script whether it's a valid route or not. For that, read the
``MOCK_REQUEST_NOT_FOUND`` environment variable. :ref:`You can read more
about here. <middlewares_not_found>`

For a complete list of all environment variables that can be read from
middleware handlers, `consult this section.
<shell_scripts.html#environment-variables-for-request-handlers>`_

Conditions for Middlewares
~~~~~~~~~~~~~~~~~~~~~~~~~~

Middlewares can use conditions, such as the ones `specified in the Conditions
Reference <conditions_reference.html>`__, in order to make custom filters. Read
further to learn more.

So far we've seen that Middlewares can use the ``route_match`` configuration
parameter in order to execute Middlewares for certain routes, but that's a very
simple kind of filter. By using the "conditions" mechanism you can define more
complex kinds of filters. For example, following is a Middleware that is only
executed when a request is made to a route that does not exist - in order
words, we're making a custom 404 page for our API:

.. code:: json

    {
      "middlewares": [
        {
          "exec": "echo 'New response body!' > $MOCK_RESPONSE_BODY",
          "condition": {
            "type": "querystring_match",
            "key_values": {
              "foo": "bar"
            }
          }
        }
      ],
      "endpoints": []
    }

The middleware above modifies all requests that have the ``foo=bar``
querystring.

.. _middlewares_not_found:

Handling "not found" requests with Middlewares
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

`mock` serves requests based on the endpoints that you configured. When a
request is made with a route for which an endpoint hasn't been set, `mock`
responds with ``405 Method not allowed`` HTTP Status.

You may want to change that behavior and provide customised responses when
requests are made to unexisting routes. That's possible through Middlewares.

By default, Middlewares are set to be executed for all route patterns, unless
you use the ``route_match`` option and filter the route patterns that you
desire. In other words, by default Middlewares are executed for all "not_found"
requests.

You can identify whether the current request being handled is valid or not by
reading the ``MOCK_REQUEST_NOT_FOUND`` environment variable. It's a boolean
variable.

Furthermore, requests for invalid routes will result in middleware having the
``MOCK_RESPONSE_STATUS_CODE`` set to ``405``. If you want a different response
status code to be set, just change that variable's value.

With that, we can easily create custom "not found" pages. Let's look at an
example on how to accomplish this:

.. code:: sh

   $ mock serve -p 3000 \
     --middleware "sh path/to/my/middleware.sh" \
     --route foo/bar \
     --response "Hello world!"

.. code:: sh

    // path/to/my/middleware.sh

    if [ "${MOCK_REQUEST_NOT_FOUND}" = "true" ]
    then
        mock set-status 404

        printf "This page does not exist!" | mock write
    fi

Let's analyse the result. Request ``GET foo/bar``:

.. code:: sh

   $ curl localhost:3000/foo/bar
   # Prints out: Hello, world!

Nothing special there - we got the exact response string as defined the
endpoint's response - no modification was made to that response through our
middleware because it only manipulates the response if it's a Not Found request
- which it detects by reading ``$MOCK_REQUEST_NOT_FOUND``.

Let us now request to a route that does not exist:

.. code:: sh

   $ curl localhost:3000/foo/bar/2
   # Prints out: This page does not exist!

This time the response was manipulated by the middleware and the custom Not
Found page was returned successfully by the server.
