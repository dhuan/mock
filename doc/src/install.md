# Installation

## Download *mock*

The easiest installation method is simply downloading *mock*'s executable file. It's built as a standalone program with no dependencies, so just choose one of the supported OSs and start using it, or check the releases page where you'll find other download formats available:

- [Download *mock* for Linux](__DOWNLOAD_LINK_LINUX__)
- [Download *mock* for MacOS](__DOWNLOAD_LINK_MACOS__)
- [Releases](https://github.com/dhuan/mock/releases)

Or, here's a few lines of shell for quickly installing it:

```sh
$ wget -O mock.tgz __LINUX_TARBALL__
$ tar xzvf mock.tgz
$ ./mock version
```

## Install through source code

Before proceeding make sure that your system has the following programs installed:

- Golang 1.18 or later
- GNU Make
- Git

Clone *mock*'s source code:

```
$ git clone https://github.com/dhuan/mock
```

Move into the new folder where *mock* was cloned run the make command for building the program:

```
$ make build
```

If executed successfully, the *mock* executable file should've been created inside the `bin` folder from the root repository path.

```
$ ./bin/mock version
```
