Middlewares
===========

APIs created with *mock* can have middlewares, which are simply shell-scripts
that are executed at the following moments:

- Immediately after a request has reached your API, before it's processed by
  endpoint handlers - at which moment you can manipulate the request object
  before it's processed by the endpoint handler; 
- Before a response is sent to the client;

The following tasks can be performed by Middlewares:

- Execute any logic that can be executed through shell scripts;
- Manipulate requests before they're handled by endpoint handlers - for
  example, adding additional headers to the request which weren't set by the
  HTTP client;
- Manipulating responses before they're sent to HTTP clients - like adding new
  content the body response or sending out new Response Headers.

Read further to learn more about it or jump to the examples section for quickly
see middlewares in action.

Setting up Middlewares
~~~~~~~~~~~~~~~~~~~~~~

Following is a *mock* configuration file containing middlewares:

.. code:: json

   {
    "middlewares": [
        {
            "exec": "sh /path/to/middleware/script.sh",
            "type": "before_response",
            "route_match": "*"
        }
    ],
    "endpoints": [
    ]
   }

Let's go over each the fields shown in the ``middleware`` objects:

- ``exec``: A shell command that will be executed which can perform a middleware operation, such as changing request/response. Anything that's valid in the shell will be valid here, such as ``sh path/to/some/script.sh`` or ``/path/to/some/program``.
- ``type``: The available Middleware Types are: ``before_request``, ``before_response``.
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
  +  --middleware-before-request path/to/some/script.sh
  +  --route-match 'foo/bar'
  +  --middleware-before-response path/to/another/script.sh

Examples of Middlewares
~~~~~~~~~~~~~~~~~~~~~~~

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

