How-tos & Examples
==================

Delaying specific endpoints
---------------------------

Making an existing API slow can be easily accomplished combining mock's
:ref:`Base APIs <base_api>` and the :ref:`delay option.
<cmd_options_reference__delay>`

.. code:: sh

    mock serve -p 8000 --base example.com --delay 2000

You may want however to make a specific endpoint slow instead of the whole API.
This can be achieved using :ref:`middlewares <middlewares>`: 

.. code:: sh

    mock serve -p 8000 --base example.com --middleware '
    if [ "${MOCK_REQUEST_ENDPOINT}" = "some/endpoint" ]
    then
        sleep 2 # wait two seconds
    fi
    '

With that last example, our API at ``localhost:8000`` will act as a proxy to
``example.com``. All requests will be responded immediately except
``some/endpoint`` which will have a delay of 2 seconds.

An API powered by multiple languages
------------------------------------

.. code:: sh

    mock serve -p 3000 \
        --route js \
        --exec '
    node <<EOF | mock write
    console.log("Hello from Node.js!")
    EOF
    ' \
        --route python \
        --exec '
    python3 <<EOF | mock write
    print("Hello from Python!")
    EOF
    ' \
        --route php \
        --exec '
    php <<EOF | mock write
    <?php
    echo "Hello from PHP!\n";
    ?>
    EOF
    '

Let's test it:


.. code:: sh

   $ curl localhost:3000/js
   # Prints out: Hello from Node.js!
   $ curl localhost:3000/python
   # Prints out: Hello from Python!
   $ curl localhost:3000/php
   # Prints out: Hello from PHP!
