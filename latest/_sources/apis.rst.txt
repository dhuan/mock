Creating APIs
=============

The simplest endpoint configuration we can define looks like this:

.. code:: json

   {
     "endpoints": [
       {
         "route": "foo/bar",
         "method": "POST",
         "response": {
           "foo": "bar"
         }
       }
     ]
   }

A ``POST`` HTTP Request to ``/foo/bar`` will respond you with
``{"foo":"bar"}``, as can be seen in the ``response`` endpoint
configuration parameter above.

Endpoint Routes can also be set with wildcards:

.. code:: diff

    {
      "endpoints": [
        {
   +      "route": "foo/bar/*",
          "method": "POST",
          "response": {
            "foo": "bar"
          }
        }
      ]
    }

With the configuration above, requests such as ``foo/bar/anything`` and
``/foo/bar/hello/world`` will be responded by the same Endpoint.

Besides wildcards, routes can have placeholder variables as well, such
as ``foo/bar/{some_variable}``. In order to read that variable and do
something useful with it, you will need to `define shell scripts that
act as handlers for your Endpoints. <#responses-from-shell-scripts>`__

In the next sections we’ll look at other ways of setting up endpoints.

Endpoints defined through command-line parameters
-------------------------------------------------

An alternative for creating configuration file exists - endpoints can be
defined all through command-line parameters. Let’s start up mock with
two endpoints, ``hello/world`` and ``hello/world/again``:

.. code:: sh

   $ mock serve \
     --route 'hello/world' \
     --method GET \
     --response 'Hello world!' \
     --route 'hello/world/again' \
     --method POST \
     --response 'Hello world! This is another endpoint.' 

As shown above, all which can be accomplished through JSON configuration
files can be done through command-line parameters, it’s just a matter of
preference. As we move forward through this manual learning more
advanced functionality, you’ll be instructed on how to achieve things in
both ways - the above only scratches the surface. A few notes to be
aware while using command line parameters:

-  Both configuration file and command-line parameters can be used
   together, but when routes are defined as parameters which have been
   already defined in the configuration file, the former will overwrite
   the latter. In other words, command-line parameters defined endpoints
   always overwrite the ones defined in config (which have the same
   route and method combination).

File-based response content
---------------------------

In the earlier example, ``response`` is a JSON object containing the
response JSON that you’ll be responded with. However, as you setup
complex APIs, your configuration file starts getting large and not
easily readable. In the following example, we’re setting the response
content by referencing a file, thus leaving the configuration file more
readable:

.. code:: diff

     {
       "endpoints": [
         {
           "route": "foo/bar",
           "method": "POST",
   +       "response": "file:path/to/some/file.json"
         }
       ]
     }

To define responses referenced by files using command-line parameters,
``--response-file`` can be used:

.. code:: diff

    $ mock serve \
      --route "foo/bar" \
      --method "POST" \
   +  --response-file path/to/some/file.json

The above can also be accomplished with
``--response "file:path/to/some/file.json"``.

.. toctree::
   :maxdepth: 2
   :caption: Contents:

   headers
   route_params
   shell_scripts
   static_files
   status_codes
   conditional_responses
   env_vars
   middlewares
   base_apis
