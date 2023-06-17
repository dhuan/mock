Reading Environment Variables
=============================

Responses can include any environment variable. The following example
starts up *mock* with a custom environment variable and includes its
variable in an endpoint’s response.

.. code:: sh

   $ FOO=BAR mock serve -c path/to/config.json

And then the configuration file:

.. code:: json

   {
     "endpoints": [
       {
         "route": "foo/bar",
         "method": "GET",
         "response": "The value of 'FOO' is ${FOO}."
       }
     ]
   }

Let’s accomplish the same but now using command-line parameters instead
of configuration file. Note here the usage of single-quotes around the
response string, because we don’t want these variables to be processed
by the shell program, but by *mock* instead:

.. code:: diff

    $ export FOO=bar

    $ mock serve \
   +  --route 'foo/bar' \
   +  --response 'The value of FOO is ${FOO}.'
