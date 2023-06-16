Route Parameters
================

Route Parameters are named route segments that can be captured as values
when defining Responses.

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

   Important: Note in the example above that the response string was
   wrapped around single-quotes, that is necessary because the variable
   ‘${book_name}’ is NOT supposed to be processed by the shell program,
   instead **mock** will process that variable while processing the
   request’s reponse, as ``book_name`` is a Route Parameter and not a
   shell variable.
