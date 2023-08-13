Base APIs
=========

A `mock` server can be configured to have a Base API, which results in all
requests received to be forwaded to the API that it's set to be based from.
When an endpoint route was configured for the requesting route, it takes
priority over the Base API. Let's see a simple example:

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
endpoint that was not configured, then the request will be forwaded to the
domain that we've set as the Base API. In other words, a request to ``foo/bar``
will result in a proxy request to ``example.com/foo/bar``.

Intercepting responses
----------------------

Middlewares can be used to manipulate responses given by a Base API. `Check the
Middlewares documentation section for learning all about them.
<middlewares.html>`__ In fact Middlewares make no distinction of whether the
current request is meant for a `mock` endpoint or an actual Base API when
executing its handler - but on your Middleware handler you can find out whether
the current Middleware execution is for a Base API request or not by reading
the ``$MOCK_BASE_API_RESPONSE`` environment variable. The middleware below adds
a header ``Foo: bar`` to all responses proxied to a Base API:

.. code:: sh

   if [ "$MOCK_BASE_API_RESPONSE" = true ];
   then
     printf "Foo: bar" >> $MOCK_RESPONSE_HEADERS
   fi

Alternatively you can just use the ``route_match`` middleware option in order
to filter the requests which you want to manipulate, targetting the route
patterns that are meant for your Base API of choice. Then again, the middleware
documentation section covers the subject in more details.

Intercepting requests
---------------------

Requests sent to Base APIs are sent exactly as defined by the HTTP client who's
requesting to `mock` - all the headers, payload etc, everything unmodified. It's
possible however to modify the request object - any of its properties such as
its headers, payload - through Middlewares.

Set the Middleware with its type as ``on_request_to_base_api`` to achieve it.
The middleware handler script is able to modify the request by writing to the
files assigned to the following environment variables:

-  **MOCK_REQUEST_BODY**: Write to this file to modify the request payload.
   It contains the payload sent by the requesting client.
-  **MOCK_REQUEST_HEADERS**: Write to this file to modify the request headers.
   It contains the headers sent by the requesting client.

Base APIs and TLS
-----------------

The Base API option can take a simple domain (``example.com``) as its value, or
a protocol+domain combo (``https://example.com``). Read further to understand
how the different methods differ:

Domain only: todo

Protocol + domain combo: todo
