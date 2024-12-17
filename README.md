# Story Guardian

`story-guardian` is a CLI tool used for periodically downloading bloom filter files. It supports retries with a delay if the download fails
and allows flexible output file paths depending on user needs.

## Features

* Automatically downloads bloom filter files every 24 hours.
* Configurable output directory for storing the downloaded bloom filter.
* Supports retrying download with delays in case of failures.
* Supports uploading the filtered reports to the CipherOwl server.
* Customizable output path based on system type (Linux, MacOS).

## Installation

To install story-guardian from source, you must have Go (1.16 or above) installed.

1. Clone the repository:

```shell
git clone https://github.com/piplabs/story-guardian.git
```

2. Navigate into the project directory:

```shell
cd story-guardian
```

3. Install the dependencies:

```shell
go mod tidy
```

4. Build the application:

```shell
go build -o story-guardian cmd/*.go
```

5. Move the binary into your PATH or run it locally:

```shell
mv story-guardian /usr/local/bin/
```

6. Configure environment variables for the client ID and secret:
```shell
export CIPHEROWL_CLIENT_ID=...
export CIPHEROWL_CLIENT_SECRET=...
```

Now you can execute `story-guardian` as a CLI tool on your terminal.

## Usage

Once the `story-guardian` has been installed, you can invoke it by running:

```shell
story-guardian [flags]
```

### Flags

The tool allows customization of where the bloom filter is saved by specifying flags. By default, the output directory
is automatically configured depending on the user's OS (`$HOME/.story/geth/guardian` on Linux
or `$HOME/Library/Story/geth/guardian` on MacOS).

You can override this by providing your own output directory using the `-o` or `--output-dir` flag.

### Available flags:

* `-o`, `--output-dir`: The directory to store the bloom filter files. (default: OS-specific,
  e.g., `$HOME/.story/geth/guardian` for Linux)

### Examples

1. *Basic usage (use default path)*: To run the program using the default output path for your system (
   i.e., `$HOME/.story/geth/guardian` on Linux, `$HOME/Library/Story/geth/guardian` on Mac):

```shell
story-guardian
```

2. *Specifying custom output directory*: To specify a custom output directory for bloom filters, use the `-o`
   or `--output-dir` flag:

```shell
story-guardian -o /path/to/custom/directory
```

3. *Running the downloader in the background*: Since this tool is designed to run periodically, you can run it in the
   background using the following:

```shell
nohup story-guardian > story-guardian.log 2>&1 &
```

This runs the downloader in the background, redirecting the output to a log file.