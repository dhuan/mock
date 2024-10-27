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

Writes data to the HTTP Response.

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

get-route-param
---------------

.. code:: sh

   $ mock get-route-param some_route_param_name

Gets a `Route Parameter <route_params.html>`__. If the parameter doesn't
exist, nothing is printed out and `mock` exists with ``1``, otherwise the
parameter value is printed out and it exits with ``0``.
