.. _shell_script_responses:

Responses from Shell scripts
============================

You can write shell scripts that will act as “handlers” for your API’s
Requests.

.. code:: json

   {
     "endpoints": [
       {
         "route": "foo/bar",
         "response": "sh:my_shell_script.sh"
       }
     ]
   }

In the example above, any request to ``POST /foo/bar`` will result in
*mock* executing the ``my_shell_script.sh``.

Scripts can define the HTTP Response through different methods:

- Write to response files which are referenced in environment variables like ``$MOCK_RESPONSE_BODY``.
- Use command-line utilities provided by `mock`, such as :ref:`mock write. <shell_utils_write>` and :ref:`mock set-header <shell_utils_set_header>`.

A `hello world` script-based response handler looks like the following:

.. code:: sh

    echo "Hello world!" | mock write

To further customize your script handlers, you may also pass parameters,
just like you can normally pass parameters in a shell command:

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
   +      "response": "sh:my_shell_script.sh some_param another_param"
        }
      ]
    }

To define responses with shell scripts using command-line parameters,
use the following:

.. code:: diff

    $ mock serve \
   +  --route "foo/bar" \
   +  --shell-script my_shell_script.sh

Alternatively, shell commands can be set as one-liners with ``exec``
instead of ``sh``, not requiring you to create a shell script file. As
an example, the endpoint below responds with a list of files of the
current folder (``ls -la``):

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
   +      "response": "exec:ls -la | mock write"
        }
      ]
    }

You can use more advanced shell functionalities within ``exec``, such as
pipes. Let’s set an endpoint that returns the amount of files that exist
on the home folder:

.. code:: diff

    {
      "endpoints": [
        {
          "route": "foo/bar",
   +      "response": "exec:ls ~ | wc -l | mock write"
        }
      ]
    }

The same can be accomplished through command-line parameters:

.. code:: diff

    $ mock serve \
   +  --route "foo/bar" \
   +  --exec 'ls | sort | mock write'

Environment Variables for Request Handlers
------------------------------------------

As with any shell programs, environment variables can be read. Besides
environment variables provided the operating system environment, there are also
variables provided by `mock` which give you useful information about the
current request being handled.

These variables `are documented in detail in the Mock Variables
<mock_vars.html>`_ page.


Route Parameters - Reading from Shell Scripts
---------------------------------------------

Route Parameters can be read from shell scripts. Suppose an endpoint
exists as such: ``user/{user_id}``. We could then retrieve the User ID
parameter by reading the ``MOCK_ROUTE_PARAM_USER_ID`` environment
variable.

Alternatively we can retrieve the same route parameter using
:ref:`get-route-param<shell_utils_get_route_param>`:

.. code:: sh

    $ mock get-route-param user_id

Response Files that can be written to by shell scripts
------------------------------------------------------

The following environment variables enable you to manipulate the HTTP Response:

-  **MOCK_RESPONSE_BODY**: A file that can be written to in order to set the
   HTTP Response.
-  **MOCK_RESPONSE_STATUS_CODE**: A file that can be written to in order to
   define the HTTP Status Code.
-  **MOCK_RESPONSE_HEADERS**: A file that can be written to in order to define
   the HTTP Headers.

In the following example, we’ll see what a Handler looks like, which
responds with a simple ``Hello world!`` body content, a ``201`` Status
Code and a few custom HTTP Headers.

.. code:: sh

   echo Hello world! > $MOCK_RESPONSE_BODY

   cat <<EOF > $MOCK_RESPONSE_HEADERS
   Some-Header-Key: Some Header Value
   Another-Header-Key: Another Header Value
   EOF

   echo 201 > $MOCK_RESPONSE_STATUS_CODE
