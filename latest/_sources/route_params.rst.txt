Route Parameters
================

Route Parameters are named route segments that can be captured as values
when defining Responses. Let's look at how to capture these parameters and use
them for our API logic.

With static response data
-------------------------

.. code:: json

   {
     "endpoints": [
       {
         "route": "books/search/author/{author_name}/year/{year}",
         "method": "GET",
         "response": "You're searching for books written by ${author_name} in ${year}."
       }
     ]
   }

With the endpoint configuration above, a request sent to
``books/search/author/asimov/year/1980`` would result in the following
response: ``You're searching for books written by asimov in 1980.``

Besides static responses as exemplified, all kinds of responses can read
Route Parameters - a response file name can be dynamic based on a
parameter:

.. code:: json

   {
     "endpoints": [
       {
         "route": "book/{book_name}",
         "method": "GET",
         "response": "file:./books/${book_name}.txt"
       }
     ]
   }

..

   Route Parameters can also be read by Shell-Script Responses. `Read
   more about it in its own guide
   section. <shell_scripts.html#route-parameters-reading-from-shell-scripts>`__

To read route parameters through endpoints defined by command-line
parameters, the same syntax applies:

.. code:: diff

    $ mock serve \
   +  --route "book/{book_name}" \
   +  --response-file 'books/${book_name}.txt'

.. note::

   In the example above that the response string was wrapped around
   single-quotes, that is necessary because the variable ‘${book_name}’ is NOT
   supposed to be processed by the shell program, instead **mock** will process
   that variable while processing the request’s reponse, as ``book_name`` is a
   Route Parameter and not a shell variable.

With shell scripts
------------------

From shell scripts, there are two different ways through which we can capture
route parameters.

We can simply read the environment variable ``MOCK_ROUTE_PARAM_<NAME>``.

Let's spin up a `mock API` with a route that has two route parameters:

.. code:: sh

   $ mock serve -p 3000 \
       --route 'say_hi/{name}/{location}' \
       --exec 'printf "Hello! My name is %s. I live on %s." "${MOCK_ROUTE_PARAM_NAME}" "${MOCK_ROUTE_PARAM_LOCATION}" | mock write'


Now let's request that route and see the result:

.. code:: sh

    $ curl -v localhost:3000/say_hi/john_doe/earth

    Hello! My name is john_doe. I live on earth.

.. note::

   If you think the "exec" command above is too lengthy and would like to keep
   your API logic more organised, you can just put that script into a shell file and instead of ``--exec``, use ``--shell-script /path/to/some/script.sh``.

.. note::

   All environment variables provided by `mock` are uppercased. Therefore if
   your route param is named ``fooBar``, the environment variable will be
   available to your script as ``MOCK_ROUTE_FOOBAR``
