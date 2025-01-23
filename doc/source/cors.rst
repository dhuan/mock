Handling CORS
=============

The ``--cors`` flag can be used when running *mock*. It will take care
of setting up all the necessary headers in your API’s Responses to
enable browser clients to comunicate without problems:

.. code:: sh

   $ mock serve --cors -c /path/to/your/config.json

Customizing OPTIONS Requests
~~~~~~~~~~~~~~~~~~~~~~~~~~~~

The ``--cors`` flag adds the necessary headers for you for all requests,
including OPTIONS requests. There may be cases however when these headers do
not satisfy your needs, and you need to further customise them. You may use
:ref:`Middlewares <middlewares>` to customize the headers for OPTIONS requests:

.. code:: sh

    mock serve \
        -p 3000 \
        --cors \
        --middleware 'test $MOCK_REQUEST_METHOD = "options" && mock set-header foo bar' \
        --route foo/bar \
        --response "Hello, world"

The middleware above verifies whether the request being handled is an OPTIONS
request, and then if it is, it adds the ``Foo: bar`` HTTP Header to the
Response.

.. note::

   For an example of when it'd be necessary to manipulate headers for OPTIONS
   requests: there are cases when it's not desireable that
   ``Access-Control-Allow-Origin`` be set to ``*``, and that it be set to an
   actual origin value.
