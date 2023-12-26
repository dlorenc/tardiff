# tardiff

This program compares two tarballs and reports any differences in their contents.

## Usage

```bash
$ tardiff tarball1 tarball2
```

## Output

The program will output any differences it finds between the two tarballs, including:

* Files that are only present in one tarball
* Files that have different names in the two tarballs
* Files that have different sizes in the two tarballs
* Files that have different PAX headers in the two tarballs
* Files that have different contents in the two tarballs

## Examples

To compare two tarballs named `tarball1.tar` and `tarball2.tar`, you would run the following command:

```bash
$ tardiff tarball1.tar tarball2.tar
```

## Installation

```bash
$ go install github.com/dlorenc/tardiff@latest
```

## License

This program is licensed under the Apache 2.0 License.
