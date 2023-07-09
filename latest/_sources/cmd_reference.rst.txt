Command-line Options Reference
==============================

Use the table of contents to navigate and see all available command-line
options.

Standard options
~~~~~~~~~~~~~~~~

``-c`` or ``--config``
----------------------

Path to a configuration file with specs of an API to be run. This is optional,
as you may instead provide an API's specs through command-line arguments.

``-p`` or ``--port``
--------------------

Defines the port through which `mock` will serve your API. Defaults to ``3000``
if not provided.

Options for specifying APIs
~~~~~~~~~~~~~~~~~~~~~~~~~~~

APIs can be specified through configuration files or command-line options. This
section cover all the available command-line options for customizing an API.
See the example below:

.. code:: sh

   $ mock serve \
     --route 'say_hi/{name}' \
     --method GET \
     --response 'Hello world! My name is ${name}.' \
     --route "what_time_is_it" \
     --method GET \
     --exec 'printf "Now it is %s" $(date +"%H:%M") > $MOCK_RESPONSE_BODY'

Find below all the available options that enable you to specify an API:

``--route``
-----------

Starts defining a new endpoint. The value for this option is the endpoint's
route.

``--method``
------------

Sets the HTTP Method.

``--response``
--------------

Sets the HTTP Response Body. Receives a string as value which will be resulting
response.

``--response-file``
-------------------

Sets the HTTP Response Body to be the contents of a file.

``--response-sh`` or ``--shell-script``
---------------------------------------

Sets a shell script file to be executed for defining the HTTP Response Body.
The value must be a path to a shell script file such as
``/path/to/some/script.sh``.

``--exec``, ``--response-exec``
-------------------------------

Sets a shell command to be executed for generating the HTTP Response Body.

``--file-server``, ``--response-file-server``
---------------------------------------------

Sets the HTTP Response to be a server of static files from a given directory.
The value must be a directory location, either relative or absolute. `Check
here for more details. <static_files.html>`__

``--header``
------------

Sets a new HTTP Response Header. The value must be formatted as such:
``Some-Header-Key: Some header value``.

``--status-code``
-----------------

Sets the HTTP Status Code. Only numbers are valid for this option.

Options for Middlewares
~~~~~~~~~~~~~~~~~~~~~~~

`Jump to the Middlewares documentation page to learn about it.
<middlewares.html>`__ Below is an example of running `mock` using Middlewares
defined through command-line parameters.

.. code:: sh

  $ mock serve \
    --route foo/bar
    --response "Hello world!"
    --middleware "sh path/to/some/script.sh"
    --route-match 'foo/bar'
    --middleware "sh path/to/another/script.sh"

``--middleware``
----------------

Sets a new Middleware. The value must be a shell command that will be executed
acting as the Middleware Handler.

``--route-match``
-----------------

Sets a regular expression for matching against routes for filtering a
Middleware. You may define Middlewares without using this option, thus setting
the Middleware to be executed for all requests.

Miscellaneous
~~~~~~~~~~~~~

``--cors``
----------

With ``--cors`` all HTTP Responses will include the necessary headers so
that your browser does not complain about cross-origin requests.

::

   $ mock serve -c /path/to/config.json --cors

``-d`` or ``--delay``
---------------------

Sets the amount of milliseconds that each request will wait before
receiving a response. When not set, requests receive responses
immediately.

The following example configures *mock* to delay every request to 3
seconds:

::

   $ mock serve -c /path/to/config.json --delay 3000
