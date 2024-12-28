.. _shell_utils:

Shell Utilities
===============

`mock` provides a set of "utilities" for manipulating request and response data
easily.

For a simple example, suppose you wanted to make a simple "search and replace"
operation before allowing the HTTP Response to be delivered to the requesting
client. We could just use `sed`:

.. code:: sh

    $ sed -i 's/foo/bar/g' $MOCK_RESPONSE_BODY

Or you could just use `mock`'s ``replace`` utility:

.. code:: sh

   $ mock replace foo bar

This documentation is a reference to all available such utilities for
manipulating requests and responses.

.. note::

    Shell utilities are also usable in :ref:`Middlewares <middlewares>`.
    Let's accomplish something similar as above, replacing text within the
    response body, but this time with :ref:`middlewares <middlewares>` instead:

    .. code:: sh

       $ mock serve -p 3000 \
           --middleware 'mock replace world universe' \
           --route foo/bar \
           --response 'Hello, world!'

    Let's request our `mock API` now and find out the result:

    .. code:: sh

       $ curl localhost:3000/foo/bar

       // Prints out: Hello, universe!

.. _shell_utils_write:

write
-----

.. code:: sh

   $ printf "Hello world!" | mock write

Writes data to the HTTP Response. Unless ``--append`` is used, ``write`` will
overwrite any response body previously defined.

Options:

- ``-a, --append``: Appends data to the response body instead of overwriting
  it.

replace
-------

.. code:: sh

   $ mock replace foo bar

Searches for `<ARG 1>` and replaces it with `<ARG 2>`, on the HTTP Response Body.

Options:

- ``--regex``: Treats *<ARG 1>* (the search parameter) as a regular expression.

wipe-headers
------------

.. code:: sh

   $ mock wipe-headers some-header-key another-header-key
   $ mock wipe-headers --regex some-pattern another-pattern

Removes one or more HTTP Headers. The header names passed as parameters must be
the exact header name. The string matching is case-insensitive.

Options:

- ``--regex``: The strings passed will be used as regex patterns for matching
  against the header keys.

.. _shell_utils_set_header:

set-header
----------

.. code:: sh

   $ mock set-header foo bar

Adds an HTTP Header to the Response. If the provided header name was already
set previously, then the provided header value will just overwrite the
previous one.

.. _shell_utils_set_status:

set-status
----------

.. code:: sh

   $ mock set-status 400

Sets the HTTP Status Code for the current response being handled.

.. _shell_utils_get_route_param:

get-route-param
---------------

.. code:: sh

   $ mock get-route-param some_route_param_name

Gets a `Route Parameter <route_params.html>`__. If the parameter doesn't
exist, nothing is printed out and `mock` exists with ``1``, otherwise the
parameter value is printed out and it exits with ``0``.

get-query
---------

.. code:: sh

   $ mock get-query
   # foo=bar&someKey=someValue
   $ mock get-query foo
   # bar

Gets a querystring value from the Request URL.

If no parameter is passed, then the whole querystring string is printed out. If
a parameter is passed then the querystring with that key is printed out. Exit
status code is 0 if a valid key is provided, otherwise 1 is returned.

If the current request being handled does not contain any querystring,
``get-query`` will print nothing, returning with status code 1.

get-header
----------

.. code:: sh

   $ mock get-header
   # Prints out all headers
   $ mock get-header authorization
   # authorization: Bearer xxx
   $ mock get-header --regex auth
   # authorization: Bearer xxx
   $ mock get-header -v authorization
   # Bearer xxx

Gets the HTTP Headers from the current request, based on your search criterias.
If no search string is passed, all headers are printed out. The search is case
insensitive. Unless ``--regex`` is used, the search string will only match if
the it's typed the full header key name.

Options:

- ``--regex``: Use regular expression for searching.
- ``-v, --value``: Print out only the header value, otherwise the whole header
  line is printed.

Exit code: If no headers are found given the search criteria, `1` is returned,
otherwise `0` when headers are found.

get-payload
-----------

.. code:: sh

   $ mock get-payload
   # Prints out all request payload
   $ mock get-payload someFieldName
   # Prints out the "someFieldName" field from the JSON payload, or multipart
   # form data, or URL-encoded data depending on the content type.

Gets either the whole request payload, or a "field" from it.

To get the whole payload, just execute ``get-payload`` without any parameter.

If a parameter is passed, `mock` will identify the kind of payload the request
is formatted in, through the `Content-Type` header, and then will extract the
field based on the name passed as parameter.

When extracting fields from the payload, the following payload formats are
supported:

- JSON payloads (``Content-Type: application/json``)
- `Multipart Form Data payloads <https://en.wikipedia.org/wiki/MIME#form-data>`_ (``Content-Type: multipart/form-data``)
- `HTML Form Data <https://en.wikipedia.org/wiki/Percent-encoding#The_application.2Fx-www-form-urlencoded_type>`_ (``Content-Type: application/x-www-form-urlencoded``)

For JSON payloads, it's possible to get nested values, with the following
syntax:

.. code:: sh

   $ mock get-payload user.name
   $ mock get-payload users[0].name

.. important::

   Retrieving nested values as exemplified above will work only with JSON
   payloads.

About the exit code:

- If no parameters are provided, the exit code will always be ``0``.
- If a parameter is provided to extract a payload field, ``0`` is returned if
  the field exists, otherwise ``1`` is returned.
