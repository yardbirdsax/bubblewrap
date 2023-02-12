# bubblewrap

`bubblewrap` is a library that wraps around [the Charmbracelet `bubbletea`
library](https://github.com/charmbracelet/bubbletea) to provide a simpler API for building CLI
applications that require input from users, without the complexity of managing the model, view,
update constructs underlying the [ELM architecture](https://guide.elm-lang.org/architecture/) that
Bubbletea and other Charmbracelet libraries are built on.

As a CLI developer, I commonly want to prompt users for things or let them choose from a menu of
options. In some cases, these CLIs need fancy / robust state management and presentation details,
which the Charmbracelet ecosystem tools are perfect for. However, in an equal (if not greater)
number of cases, I want the elegance and flexibility of the Charmbracelet tools, but don't need to
do things like store / mutate / view the state of things done by the CLI. I just want present users
with a nice looking, consistently styled interface to get the information I need, then do the work
in simpler ways without fancy TUIs. This is why I wrote `bubblewrap`, so that I could have an easy
API / entrypoint into common operations without constantly writing the same boilerplate code around
things like handling different kinds of key presses.

Currently `bubblewrap` contains a very simple abstraction of [the `textinput`
bubble](https://github.com/charmbracelet/bubbles/#text-input). Over time and as needs arise I might
add additional wraps for things like choosing items in a list or multiline text. PRs for additional
features are welcome!

Examples can be found in the [`examples` directory](examples/).
