mock - Language-agnostic API mocking and testing utility
========================================================

*mock* enables you to quickly set up HTTP servers for end-to-end tests.

-  Define endpoints and their respective responses through easy syntax;
-  Make assertions on…

   -  Whether a given endpoint was requested;
   -  If a JSON payload body was passed correctly to a given endpoint;
   -  If a header value was passed correctly;
   -  And other useful things…

.. code:: sh

   $ mock serve --port 3000 \
     --route 'say_hi/{name}' \
     --method GET \
     --response 'Hello world! My name is ${name}.' \
     --route "what_time_is_it" \
     --method GET \
     --exec 'printf "Now it is %s" $(date +"%H:%M") > $MOCK_RESPONSE_BODY'

Run the example command the above and try these URLs in your browser or
any preferred HTTP client: ``http://localhost:3000/say_hi/mock`` and
``http://localhost:3000/what_time_is_it``

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

----

*Why "language-agnostic"???* - similar tools exist out there but they somehow require you to write in some programming-language when setting up fake APIs for testing. *mock* on the other hand enables you to set things up easily by just writing configuration files.


License
-------

*mock* is licensed under MIT. For more information check the `LICENSE
file. <https://github.com/dhuan/mock/blob/master/LICENSE>`__

.. toctree::
   :hidden:
   :maxdepth: 2

   install
   changelog

.. toctree::
   :hidden:
   :maxdepth: 2
   :caption: User Guide:

   apis
   test_assertions
   api_reference
   cmd_reference
   conditions_reference

.. toctree::
   :hidden:
   :maxdepth: 2
   :caption: Miscellaneous:

   cors
