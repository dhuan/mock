Welcome to mock’s documentation!
================================

*mock* enables you to quickly set up HTTP servers for end-to-end tests.

-  Define endpoints and their respective responses through easy syntax;
-  Make assertions on…

   -  Whether a given endpoint was requested;
   -  If a JSON payload body was passed correctly to a given endpoint;
   -  If a header value was passed correctly;
   -  And other useful things…

.. code:: sh

   $ mock serve --port 3000 \
     --route 'time_in/{country}' \
     --method GET \
     --exec 'zdump ${country}' \
     --route 'whois/{domain}' \
     --method GET \
     --exec 'whois ${domain}'

Run the example command the above and try these URLs in your browser or
any preferred HTTP client: ``http://localhost:3000/time_in/Japan`` and
``http://localhost:3000/whois/google.com``

Quick links
-----------

-  `Download mock for Linux <VAR_DOWNLOAD_LINK_LINUX>`__
-  `Download mock for MacOS <VAR_DOWNLOAD_LINK_MACOS>`__
-  `Releases <https://github.com/dhuan/mock/releases>`__
-  `mock\ ’s source code <https://github.com/dhuan/mock>`__
-  `Report bugs <https://github.com/dhuan/mock/issues>`__

Read further...
---------------

The core functionalities of *mock* are documented each in their
respective sections. Read further to learn:

-  `Creating APIs <apis.html>`__
-  `Test Assertions <test_assertions.html>`__

License
-------

*mock* is licensed under MIT. For more information check the `LICENSE
file. <https://github.com/dhuan/mock/blob/master/LICENSE>`__

.. toctree::
   :maxdepth: 2

   install
   changelog

.. toctree::
   :maxdepth: 2
   :caption: User Guide:

   apis
   test_assertions
   api_reference
   cmd_reference
   conditions_reference

.. toctree::
   :maxdepth: 2
   :caption: Miscellaneous:

   cors
