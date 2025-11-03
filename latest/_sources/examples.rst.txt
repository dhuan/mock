How-tos & Examples
==================

.. note::

   If you try these examples on your computer by copying/pasting the code
   snippets from this page, be aware that it'll only work if "mock" is
   installed as a global executable.

A "Who is" Lookup Service
-------------------------

.. code:: sh

    mock serve -p 3000 \
        --route 'whois' \
        --exec 'whois $(mock get-query domain) | mock write'

Let's now test it:

.. code:: sh

    $ curl localhost:3000/whois?domain=google.com
    # Prints out:
    # Domain Name: GOOGLE.COM
    # Registry Domain ID: 2138514_DOMAIN_COM-VRSN
    # ...

An API powered by multiple languages
------------------------------------

.. code:: sh

    $ mock serve -p 3000 \
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

A stateful API
--------------

.. code:: sh

    $ export TMP=$(mktemp)
    $ printf "0" > "${TMP}"

    $ mock serve -p 3000 \
        --route '/hello' \
        --exec '
    printf "%s + 1\n" "$(cat ${TMP})" | bc | sponge "${TMP}"

    printf "This server has received %s request(s) so far." "$(cat '"${TMP}"')" | mock write
    '

Let's test it:

.. code:: sh

   $ curl localhost:3000/hello
   # Prints out: This server has received 1 request(s) so far.
   $ curl localhost:3000/hello
   # Prints out: This server has received 2 request(s) so far.
   $ curl localhost:3000/hello
   # Prints out: This server has received 3 request(s) so far.

A CRUD API with a simple data storage
-------------------------------------

The following API does two tasks: add users and fetch users.

.. warning::

   To run this example, you'll need to have `jq <https://github.com/jqlang/jq>`_ installed, as the "/users"
   endpoint uses it for parsing.

.. code:: sh

    $ export DATA_DIR=$(mktemp -d)

    $ mock serve -p 3000 \
        --route 'user' \
        --method POST \
        --exec '
    # Insert a new user

    USER_NAME=$(mock get-payload name)
    USER_EMAIL=$(mock get-payload email)
    NEW_USER_ID="$(ls $DATA_DIR | wc -l | sed "s/$/+1/" | bc)"

    printf "New user ID generated: %s\n" "${NEW_USER_ID}"

    printf '"'"'{"name":"%s","email":"%s"}'"'"' "${USER_NAME}" "${USER_EMAIL}" > "${DATA_DIR}/${NEW_USER_ID}.json"

    printf '"'"'{"id":"%s"}'"'"' "${NEW_USER_ID}" | mock write
    ' \
        --route 'user/{user_id}' \
        --exec '
    # Get an existing user

    USER_ID="$(mock get-route-param user_id)"
    USER_FILE="${DATA_DIR}/${USER_ID}.json"

    if [ ! -f "${USER_FILE}" ]
    then
        mock set-status 400

        exit 0
    fi

    mock write < "${USER_FILE}"
    ' \
        --route 'users' \
        --exec '
    # Gets ALL users

    cat $DATA_DIR/*.json | jq -s | mock write
    '

Let's now test it:

.. code:: sh

    $ curl -X POST localhost:3000/user \
        -H 'Content-Type: application/json' \
        -d @- <<EOF
    {"name":"John Doe","email":"john.doe@example.com"}
    EOF
    # Prints out: {"id":"1"}

    $ curl -X POST localhost:3000/user \
        -H 'Content-Type: application/json' \
        -d @- <<EOF
    {"name":"Jane Doe","email":"jane.doe@example.com"}
    EOF
    # Prints out: {"id":"2"}

    $ curl -v localhost:3000/user/1
    # Prints out: {"name":"John Doe","email":"john.doe@example.com"}
    $ curl -v localhost:3000/user/2
    # Prints out: {"name":"Jane Doe","email":"jane.doe@example.com"}

    $ curl -v localhost:3000/user/10
    # will fail with 400/BadRequest

    $ curl -v localhost:3000/users
    # Prints out: [
    #  {"name":"John Doe","email":"john.doe@example.com"},
    #  {"name":"Jane Doe","email":"jane.doe@example.com"}
    # ]

Delaying specific endpoints
---------------------------

Making an existing API slow can be easily accomplished combining mock's
:ref:`Base APIs <base_api>` and the :ref:`delay option.
<cmd_options_reference__delay>`

.. code:: sh

    $ mock serve -p 8000 --base example.com --delay 2000

You may want however to make a specific endpoint slow instead of the whole API.
This can be achieved using :ref:`middlewares <middlewares>`: 

.. code:: sh

    $ mock serve -p 8000 --base example.com --middleware '
    if [ "${MOCK_REQUEST_ENDPOINT}" = "some/endpoint" ]
    then
        sleep 2 # wait two seconds
    fi
    '

With that last example, our API at ``localhost:8000`` will act as a proxy to
``example.com``. All requests will be responded immediately except
``some/endpoint`` which will have a delay of 2 seconds.

