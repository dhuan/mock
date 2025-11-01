.. _base_api:

Base APIs
=========

A `mock` server can be configured to have a Base API. All requests sent to the
server will be forwarded to the Base API:

.. code:: sh

   $ mock serve \
     --base "example.com" \
     --route 'hello/world' \
     --response 'Hello world!'

Alternatively you may use the ``base`` configuration option for achieving the
same:

.. code:: json

   {
     "base": "example.com",
     "middlewares": [],
     "endpoints": [
        {
          "route": "hello/world",
          "response": "Hello world!"
        }
     ]
   }

Above we've set up an API that uses ``example.com`` as its Base API,
furthermore an endpoint was set routed as ``hello/world``. If a request is made
to ``GET hello/world``, the server will act normally responding with the
response which was set - ``Hello world!``. However if a request is made to an
endpoint that has not been set, then the request will be forwaded to the Base
API. In other words, a request to ``foo/bar`` will result in a proxy request to
``example.com/foo/bar``.

.. note::

   Base APIs defined through the command-line flag takes precedence over one
   set through configuration file. There can only be one Base API for a running
   `mock` instance.

.. note::

   You can use `mock` solely as a Base API. Don't define any endpoints and
   start up `mock` with only the Base API option. It will act as a proxy to
   some other service.

Intercepting responses
----------------------

Endpoints defined for a `mock API` overwrite the ones from a Base API. In other
words, when both your API *and* Base API offer the same endpoint, then your
API's endpoint takes precedence, therefore the Base API's endpoint does not
even get requested in that scenario. But that is the default behavior only, we
can set things to behave differently if desired.

`Through response shell scripts <shell_scripts.html>`__ we can forward the
current request to the Base API, and then the HTTP Response from Base API will
be available in the ``$MOCK_RESPONSE_BODY`` environment variable, enabling us
to tweak the response if desired, before sending to the client. This is
accomplished using the `forward` command, executed from your endpoint's shell
script:

.. code:: sh

    $ mock forward

.. note::

   If the HTTP Response obtained from `forward` is encoded with gzip, `mock`
   decodes it for you, therefore the ``$MOCK_RESPONSE_BODY`` file will contain
   the decoded data for you to manipulate as you wish. Note also that `mock`
   removes the `Content-Encoding` HTTP Header upon decoding the data, therefore
   if you wish to respond to the client with encoded data, you must manually
   add again the ``Content-Encoding`` header.

Intercepting responses through middlewares
------------------------------------------

:ref:`Middlewares <middlewares>` can be used to manipulate responses given by a
Base API. In fact Middlewares make no distinction between requests to Base API
or otherwise.

On your middleware handler you can find out whether the context is that of a
Base API request or not by reading the ``$MOCK_BASE_API_RESPONSE`` environment
variable.

The middleware below adds a header ``Foo: bar`` to all responses proxied to a
Base API:

.. code:: sh

   if [ "$MOCK_BASE_API_RESPONSE" = true ];
   then
     printf "Foo: bar" >> $MOCK_RESPONSE_HEADERS
   fi

Alternatively you can just use the ``route_match`` middleware option in order
to filter the requests which you want to manipulate, targetting the route
patterns that are meant for your Base API.

Base APIs and TLS
-----------------

The Base API option can take a simple domain (``example.com``) as its value, or
a protocol+domain combo (``https://example.com``). Read further to understand
how the different methods differ:

Domain only: The protocol set by the requesting client will be respected. If a
client requests `mock` using HTTPS, then `mock` will request the Base API using
HTTPS as well.

Protocol + domain combination: If a protocol is set in the Base API's value,
then `mock` will always use that protocol when forwarding the request,
independent of the protocol chosen by requesting client.

Manipulating headers
--------------------

There may be cases when you don't want certain HTTP Headers from the Base API.
For that, `mock` provides a command for easily removing headers:

.. code:: sh

   $ mock forward
   $ mock wipe-headers some-header-key another-header-key

.. note::

    `wipe-headers` is meant to be used `inside shell scripts. <shell_scripts.html>`__

A response handler shell script using `wipe-headers` as exemplified above will
remove HTTP Headers `some-header-key` and `another-header-key`.

.. note::

    `wipe-headers` is just a faster way of manipulating the
    `$MOCK_RESPONSE_HEADERS` file. The exact same could've been accomplished
    with:

    .. code:: sh

        $ mock forward
        $ grep -v \
            -e some-header-key \
            -e another-header-key \
            $MOCK_RESPONSE_HEADERS | sponge $MOCK_RESPONSE_HEADERS

Regular expressions are also supported:

.. code:: sh

   $ mock wipe-headers --regex some-regex-pattern another-regex-pattern
