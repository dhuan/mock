Mock Variables
==============

When defining responses for endpoints, you're provided with a set of variables
containing useful information about the current request being handled. The
syntax for printing out a variable is ``${VARIABLE_NAME}``. Following is a
simple example of an endpoint that prints out the request's querystring:

.. code:: sh

    $ mock serve \
        --route foo/bar \
        --response 'The querystring: ${MOCK_REQUEST_QUERYSTRING}'

.. note::

   Note the usage of single quotes in the example above when defining the
   response! When using `mock variables` single quotes are necessary otherwise
   your shell will try to process these variables. The variables are supposed
   to be processed by `mock`, therefore single quotes must be used to avoid
   being replaced by your shell program.

Besides reading variables in responses as exemplified above, other places where
variables can be used include:

- `Response Shell Script Handlers <shell_scripts.html>`_
- `Middleware Script Handlers <middlewares.html>`_


Variable Reference
------------------

Find below all available variables.


MOCK_HOST
~~~~~~~~~

The hostname + port combination to which Mock is currently listening. (ex:
``localhost:3000``)

MOCK_REQUEST_URL
~~~~~~~~~~~~~~~~

The full URL. (ex: ``http://localhost/foo/bar``)

MOCK_REQUEST_ENDPOINT
~~~~~~~~~~~~~~~~~~~~~

The endpoint extracted from the URL. (ex: ``foo/bar``)

MOCK_REQUEST_HOST
~~~~~~~~~~~~~~~~~

The hostname + port combination that the request was sent to. (ex:
``example.com:3000``)

MOCK_REQUEST_HEADERS
~~~~~~~~~~~~~~~~~~~~

A file path containing all HTTP Headers.

MOCK_REQUEST_HEADER_FOOBAR
~~~~~~~~~~~~~~~~~~~~~~~~~~

A variable holding an individual header value. For example, if a request is
received with the header key/value as ``Foo: bar``, then this header value can
be obtained by reading the ``MOCK_REQUEST_HEADER_FOO`` environment variable.

Note that since environment variables cannot have dash characters (``-``),
`mock` converts them to underscore (``_``), for example, a header key named
`Some-header` is readable as ``MOCK_REQUEST_HEADER_SOME_HEADER``.

MOCK_REQUEST_BODY
~~~~~~~~~~~~~~~~~

For `Response Script Handlers <shell_scripts.html>`_, this variable is a file
path, containing the Request’s Body. For static responses (such as JSON or
plain text), this variable holds the actual request payload string.

MOCK_REQUEST_QUERYSTRING
~~~~~~~~~~~~~~~~~~~~~~~~

The Request’s Querystring if it exists.

MOCK_REQUEST_QUERYSTRING_KEY_NAME
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

A variable holding an individual querystring value named `KEY_NAME`. If a
request was made with the ``?foo=bar``, you can capture the "foo" parameter by
reading the variable ``MOCK_REQUEST_QUERYSTRING_FOO``.

MOCK_REQUEST_METHOD
~~~~~~~~~~~~~~~~~~~

A string indicating the Request’s Method.

MOCK_REQUEST_NTH
~~~~~~~~~~~~~~~~

A number indicating Request’s position in the request history. For example, if
two requests have been made to the ``foo/bar`` endpoint ever since *mock*
started, this being the 2nd request, the number in this variable will be 2.

MOCK_REQUEST_HTTPS
~~~~~~~~~~~~~~~~~~

This is set to `true` in case the receiving request is using HTTPS instead of
HTTP.
