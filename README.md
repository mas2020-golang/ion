# ion <!-- omit in toc -->
Ion is an all-in-one application to sum up a lot of useful tools in a single command. The swiss knife for every SysAdmin/DevOps!

## Table of Content <!-- omit in toc -->
- [Principles](#principles)
- [Getting started (Mac/Linux)](#getting-started-maclinux)
  - [Install via homebrew (Mac/Linux)](#install-via-homebrew-maclinux)
  - [Manual installation (Windows)](#manual-installation-windows)
- [File commands](#file-commands)
  - [Search command](#search-command)
- [Security commands](#security-commands)
- [Contribution on Ion](#contribution-on-ion)

## Principles

The principle of `ion` is to be light, simple and easy to use. The documentation is fully contained in the application help
that can be read as with other modern CLI tools such as docker, kubernetes, hugo. Simply type:

```shell
ion --help
```

or for a specific command help type:

```shell
ion <command> --help
```

> Note
Some commands, for a more comfortable use, have the same name as the corresponding linux/unix command. It doesn't mean that they have
the same options or the same complexity. `ion` was born to be light and easy to use. For specific and complex use cases take a look at the
Linux/Unix commands.

## Getting started (Mac/Linux)

You can install `ion` using the installing script:

```shell
curl -sS https://raw.githubusercontent.com/mas2020-golang/ion/main/install.sh | bash
```

### Install via homebrew (Mac/Linux)

To install with `homebrew` (on MacOS and Linux) first install `homebrew` package manager itself: take a look at the [official site](https://brew.sh/).

To install the application type:
```shell
brew tap mas2020-golang/ion
brew install ion
```

### Manual installation (Windows)

As of today, there doesn't exist an automatic way to install the software on Windows. It is compatible, but you have to do it manually. Go [here](https://github.com/mas2020-golang/ion/releases) to take the latest version, download the Windows file that is compliant with your CPU/OS.

Suggestion: once you have downloaded the file, you have to extract and place it in a folder that is automatically checked out by Windows. By this way you can run directly as:
```shell
ion.exe version
```
We are planning to integrate with a specific package manager also for Windows (e.g. `chocolatey`).

## File commands

Follow the list of all the *file* available commands in the current version of `ion`:
- `tail`: to show the rows of a file/standard input starting from the bottom
- `tree`: to show folders and files in a graphical representation
- `count`: to count words and lines of a specific file/standard input
- `search`: to search a single pattern into the given file/standard input
- `rm`: removes the files or folders given as an input

### Search command
The `search` command can read directly from the standard input, one or more files or directories passed an argument:

```shell
ion search [FLAGS] <PATTERN> <PATH> [ …]
```

The PATH can be one or more files and folders separated by space. The command searches for the regexp PATTERN contained in the input.
The search command returns the line in match with the pattern. In case more files are given, the output is grouped by each file and
the matched pattern is highlighted.

**Flags** are:
- **--no-colors**: no highlight colors in the output
- **--insensitive**: the search is case insensitive
- **--after <NUMBER>**: shows also the NUMBER of lines after the match
- **--before <NUMBER>**: shows also the NUMBER of lines before the match 
- **--recursive**: if the PATH is a folder searches in the sub folders too and not only in its first level
- **--invert**: shows the lines that don't match with the pattern. With this flag the count-lines, count-pattern and only-match flags are disabled.
- **--count-lines**: shows only how many lines match with the pattern (output comes first then count-pattern in case are both present)
- **--count-pattern**: shows only how many time a pattern is in match
- **--only-match**: shows only the substring that match, not the entire line
- **--only-filename**: shows only the filename when a pattern matches one or several times
- **--only-result**: if there is at least one match it returns 1, otherwise 0

**Features**

- The flags `--before` and `--after` do not accept negative numbers, in those cases the values are ignored.

- The flags `--words` for the exact correspondence with the pattern has not been implemented (as in `grep` for example). You
can reach out the same result with the regexp. For example, suppose you are interested in the `app` only and you have this
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

- `ion` doesn't accept multiple patterns, having a regex as a search engine you can get the same result with the '|' operator. When this is not possible you have to search multiple times.
- when the same pattern is searched on more files, the file path will be showed onto the standard output in this way:

  ```shell
  $ ion search 'echo' Makefile Makefile-test
  > on 'Makefile':
          @echo "==> ion test..."
          @echo "start building..."
  > on 'Makefile-test':
          @echo "==> ion test..."
  ```

- when the **PATH is a folder** ion searches only into the first level of it, unless the `--recursive` is given. You can have multiple PATHs as input. The PATH can be a single file or a folder. Example:

  ```shell
  $ ion search 'line2' /Users/andrea/Downloads/tmp/search --recursive
  > '/Users/andrea/Downloads/tmp/search/file2':
  line2
  > '/Users/andrea/Downloads/tmp/search/level2/file2':
  line2
  ```

## Security commands

Follow the list of all the *security* available commands in the current version of `ion`:
- `decrypt`: an easy way to decrypt file and folders
- `encrypt`: an easy way to encrypt file and folders

## Contribution on Ion

We really appreciate any contribution to the software. Pls, follow the [CONTRIBUTING](CONTRIBUTING) file that you find at the root level of the project.
