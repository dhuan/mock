Command-line Options Reference
==============================

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
