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
- --insensitive: the search is case insensitive
- --after <NUMBER>: shows also the NUMBER of lines after the match
- --before <NUMBER>: shows also the NUMBER of lines before the match 
- --recursive: if the PATH is a folder searches in the sub folders too and not only in its first level
- --invert: shows the lines that doesn’t match with the pattern. With this flag the count-lines, count-pattern and only-match flags are disabled.
- --count-lines: shows only how many lines match with the pattern (output comes first then count-pattern in case are both present)
- --count-pattern: shows only how many time a pattern is in match
- --only-match: shows only the substring that match, not the entire line
- --only-filename: shows only the filename when a pattern matches one or several times
- --only-result: if there is at least one match it returns 1, otherwise 0

**Features**

- The flags `--before` and `--after` do not accept negative numbers, in those cases the values are ignored.

- The flags `--words` for the exact correpondence with the patten has not been implemented (as in `grep` for example) because
you can reach out the same result with the regexp. For example, suppose you are interested in the `app` only and you have this
text:

  ```shell
  --only-result: if there is at least one match it returns 1, otherwise 0
  app1username=app1login app
  app1password=S0methingS@Str0ng!
  ```

  Using `grep` you type:

  ```shell
  grep -w 'app' test/test-files/search.txt
  ```

  Using `ion` you can reach the same result as:

  ```shell
  ion search ' app| app |app ' test/test-files/search.txt
  ```

- `ion` doesn't accept multiple patterns, having a regex as search engine you can get the same result with the '|' operator. In case this is not possible you have to search more times.
- when the same pattern is searched on more that one file the file path will be showed onto the standard output in this way:

  ```shell
  $ ion search 'echo' Makefile Makefile-test
  > on 'Makefile':
          @echo "==> ion test..."
          @echo "start building..."
  > on 'Makefile-test':
          @echo "==> ion test..."
  ```

- when the **PATH is a folder** ion searches in the first level unless the `--recursive` is given. You can have multiple PATHs as input. The PATH can be a single file or a folder. Example:

  ```shell
  $ ion search 'line2' /Users/andrea/Downloads/tmp/search --recursive
  > '/Users/andrea/Downloads/tmp/search/file2':
  line2
  > '/Users/andrea/Downloads/tmp/search/level2/file2':
  line2
  ```
