# ghibp

ghibp is small command-line utility that uses the HaveIbeenPwned database to
return information on leaked passwords or breached sites.

Powered by https://haveibeenpwned.com

## Documentation

Markdown documentation for available commands can be found under [docs/](docs/).
Additionally as part of the build manpages are generated under
[dist/man/](dist/man/).

## Build & Install

To build and install the binary, markdown docs and manpages simply issue:

```
make
sudo make install # If installing under the default PREFIX=/usr/local/share
```

The `build`, `man` and `docs` targets each build the corresponding artifacts
only. To clean up a build use `make clean` and to uninstall
`sudo make uninstall` for the default prefix.

To install/uninstall under a different prefix make use of the `PREFIX` variable.
