# ion
Ion is a all-in-one application to sum up a lot of useful tools in a single command. The swiss knife for every SysAdmin/DevOps!

## Principles

The principles of `ion` is to be light, simple and easy to use. The documentation is contained in the application itself
that be can read as you do with other modern CLI tools as docker, kubernetes, hugo. Simply type:

```shell
ion --help
```

or for a specific command help type:

```shell
ion <command> --help
```

> Note
Some commands, for a more comfortable use have the same name of the corresponding linux/unix command. It doesn't mean that have
the same options or the same complexity. `ion` is born to be light and easy to use, for specific and complex use cases take a look at the
Linux/Unix commands.

## Getting started

You can install `ion` using the installing script for Mac and Linux:

```shell
#TODO: continue from here, take note how to download the install.sh and execute...
```

## File commands

Follow the list of all the file available commands in the current version of `ion`:
- `tail`: to show the rows of a file/standard input starting from the bottom
- `tree`: to show folders and files in a graphical representation
- `count`: to count words and lines of a specific file/standard input
- `search`: to search a single pattern into the given file/standard input

### Search command
The command to exec a search in the standard input is search:

```shell
ion search [FLAGS] <PATTERN> <PATH> [ …]
```

At the moment (v0.3.0) the search command cannot search into a folder, this feature will be implemented soon.
The command searches for the PATTERN in the PATH (can be one or more paths). The PATTERN is a regular expression.
The output is done by the lines of the input that contains the pattern. In case the search is run on more files the output is grouped by each file. The matched pattern is highlighted.

**Flags** are:
- --no-colors: no highlight colors in the output
- --verbose: gives back the number of found occurrences and some details on the PATH
- --insensitive: the search is case insensitive
- --words: search only for an entire word matching
- --after <NUMBER>: shows also the NUMBER of lines after the match
- --before <NUMBER>: shows also the NUMBER of lines before the match 
- --recursive: when the input is a folder searches in the sub folders too
- --invert: shows the lines that doesn’t match with the pattern. With this flag the count-lines, count-pattern and only-match flags are disabled.
- --count-lines: shows only how many lines match with the pattern (output comes first then count-pattern in case are both present)
- --count-pattern: shows only how many time a pattern is in match
- --only-match: shows only the substring that match, not the entire line
- --only-filename: shows only the filename when a match is present one or more times
- --only-result: if there is at least one match it returns 1, otherwise 0
