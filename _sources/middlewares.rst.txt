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

Examples of Middlewares
~~~~~~~~~~~~~~~~~~~~~~~

.. note::

   The examples below are provided as shell scripts programs. But remember that
   middlewares are not limited to being shell scripts only. Any executable
   program at all can be used as a middleware.

Modify Response Body
--------------------

The following Middleware Script replaces all occurrences of ``foo`` with ``bar``
in the response body:

.. code:: sh

   sed -i 's/foo/bar/g' $MOCK_RESPONSE_BODY

Adding new headers before sending response to client
----------------------------------------------------

The following Middleware Script adds two header fields to the response:

.. code:: sh

   echo 'Header-One: Value for header one' >> $MOCK_RESPONSE_HEADERS
   echo 'Header-Two: Value for header two' >> $MOCK_RESPONSE_HEADERS

.. note::

   Observe that the script is *appeding* new headers to the headers file
   instead of overwriting it, which is indicated by the ">>" shell operator
   which appends to file. If you overwrite the file with ">", all headers
   previously set by the response handler will be overwritten.

Removing headers from the response
----------------------------------

The following Middleware Script removes all Headers that have the word
``foobar`` in their names or values.

.. code:: sh

   TMP=$(mktemp)
   cat $MOCK_RESPONSE_HEADERS | grep -v foobar > $TMP
   cat $TMP > $MOCK_RESPONSE_HEADERS


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
