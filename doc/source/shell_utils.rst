.. _shell_utils:

Shell Utilities
===============

`mock` provides a set of "utilities" for manipulating response data easily.
Read further to learn about them.

`In the previous section <shell_scripts.html>`__ we've seen that we can easily
create shell scripts that act as response handlers for API endpoints. Defining
the response data is done through writing to certain files such as
``$MOCK_RESPONSE_BODY``. Although writing to these files is easy enough, for
more complex requirements we need to write complex shell commands to accomplish
things - for example replacing strings can be achieved using `sed`:

.. code:: sh

   $ sed 's/foo/bar/g' $MOCK_RESPONSE_BODY | sponge $MOCK_RESPONSE_BODY

Although the above is simple enough, you may not be knowledgeable about all
shell tricks, not to mention that different UNIX environments may have
inconsistent or incompatible tools (GNU sed is not exactly totally compatible
with BSD sed etc).

Using `mock`'s shell utilities can save you of that burden. Let's accomplish
the same http response modification using just `mock` instead of `sed`:

.. code:: sh

   $ mock replace foo bar

Note how we didn't need to bother typing the file path as before.

In the following sections we'll look at each such utility.

.. note::

    You can also use all these utilities in :ref:`Middlewares <middlewares>`.
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
   # Prints out the "someFieldName" from the JSON request payload

Gets the request payload. If no parameters are given, the whole request
payload is printed out.

If a paramater is passed:

- If the request contains a JSON payload and JSON header, `get-payload` will
  print out the JSON field according to the provided parameter.
- If the request is a multipart/form-data one, `get-payload` will extract the
  value accordingly.

About the exit code:

- If no parameters are provided, the exit code will be always be ``0``.
- If a parameter is provided to extract a payload field, ``0`` is returned if
  the field exists, otherwise ``1`` is returned.

.. warning::

   When extracting fields from the payload, `mock` respects the content-type
   header. That means a request may contain a JSON payload, however if the
   request header is not properly set as JSON, `mock` won't give you the
   desired value.
