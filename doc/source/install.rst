Installation
============

Download *mock*
---------------

The easiest installation method is simply downloading *mock*\ ’s
executable file. It’s built as a standalone program with no
dependencies, so just choose one of the supported OSs and start using
it, or check the releases page where you’ll find other download formats
available:

-  `Download mock for Linux <VAR_DOWNLOAD_LINK_LINUX>`__
-  `Download mock for MacOS <VAR_DOWNLOAD_LINK_MACOS>`__
-  `Releases <https://github.com/dhuan/mock/releases>`__

Or, here’s a few lines of shell for quickly installing it:

.. code:: sh

   $ wget -O mock.tgz VAR_DOWNLOAD_LINK_LINUX
   $ tar xzvf mock.tgz
   $ ./mock version

Install through source code
---------------------------

Before proceeding make sure that your system has the following programs
installed:

-  Golang 1.18 or more recent
-  GNU Make
-  Git

Clone *mock*\ ’s source code:

.. code:: sh

   $ git clone https://github.com/dhuan/mock

Move into the new folder where *mock* was cloned run the make command
for building the program:

.. code:: sh

   $ make build

If executed successfully, the *mock* executable file should’ve been
created inside the ``bin`` folder from the root repository path.

.. code:: sh

   $ ./bin/mock version

Install through Go
------------------

This installation method requires only having Go installed. No need to clone the repository. Just use Go's `install` utility:

.. code:: sh

    $ go install github.com/dhuan/mock@latest

`mock` should then be available as an executable located in `$GOPATH/bin <https://go.dev/wiki/GOPATH>`_.
