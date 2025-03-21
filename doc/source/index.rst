mock - API and testing utility
==============================

.. grid:: 2

   .. grid-item::
        :columns: 7

        *mock* is an API utility - it lets you:

        - define API routes easily through API configuration files or through
          command-line parameters.
        - use shells scripts as response handlers. Or any other type of program can act
          as response handlers.
        - test your API - make assertions on whether an endpoint was requested.


   .. grid-item-card:: Quick Links
        :columns: 5

        -  `Download mock for Linux <VAR_DOWNLOAD_LINK_LINUX>`__
        -  `Download mock for MacOS <VAR_DOWNLOAD_LINK_MACOS>`__
        -  `Releases <https://github.com/dhuan/mock/releases>`__
        -  `mock\ ’s source code <https://github.com/dhuan/mock>`__
        -  `Report bugs <https://github.com/dhuan/mock/issues>`__

Let's look at a simple example - an API with 2 routes ``GET say_hi/{name}`` and
``GET what_time_is_it``:

.. code:: sh

   $ mock serve --port 3000 \
     --route 'say_hi/{name}' \
     --method GET \
     --response 'Hello, world! My name is ${name}.' \
     --route "what_time_is_it" \
     --method GET \
     --exec 'printf "Now it is %s" $(date +"%H:%M") | mock write'

Now try requesting your `mock API` at port 3000 (can also be from your
browser!):

.. code:: sh

   $ curl localhost:3000/say_hi/john_doe

   Hello, world! My name is john_doe.

   $ curl localhost:3000/what_time_is_it

   Now it is 22:00

*mock* lets you also extend other APIs (or any HTTP service, for that matter.)
Suppose you want to add a new route to an existing API running at
``example.com``:

.. code:: sh

   $ mock serve --port 3000 \
     --base example.com \
     --route 'some_new_route' \
     --method GET \
     --exec 'printf "Hello, world!" | mock write' 

With the ``--base example.com`` option above, your *mock API* will act as proxy
to that other website, and extend it with an extra route ``GET
/some_new_route``. Look up "Base APIs" in the docs for more details.

There are many other ways of further customising your APIs with *mock*. Read
further through this guide to learn.

Read further...
---------------

The core functionalities of *mock* are documented each in their
respective sections. Read further to learn:

-  :ref:`Getting started <getting_started>`
-  :ref:`Test Assertions <test_assertions>`

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
