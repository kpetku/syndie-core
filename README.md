# syndie-core
This program contains the CLI, network, database backend for syndie-gui.

Some minor utilities like a HTTP "gateway" bridge are also included.

Major WIP, not much to see here.  Expect major refactoring and breaking changes.

# Warning
This is highly experimental software that you should not use in this state.

# Dependencies
* [boltdb](https://github.com/boltdb/bolt)

# TODO
* Refactor "fetcher", "gateway", et al into Go plugins