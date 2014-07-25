# mocleaner

[![Build Status](https://drone.io/github.com/dietmarw/mocleaner/status.png)](https://drone.io/github.com/dietmarw/mocleaner/latest)

Rewrite of [Trim Trailing White Spaces](https://github.com/dietmarw/trimtrailingwhitespace) in [Go language](http://golang.org)

## Releases
You can find the pre-releases under the
[releases link](../../releases).

## Installation
Simply download the statically compiled executables for your platform from the
[releases link](../../releases).

## Usage

Currently the `mocleaner` does only remove trailing white spaces from files
of type "text/plain" (it checks for them) recursively inside `<dir>`

### Linux/Mac/FreeBSD

```
mocleaner <dir>
```
### Windows

```
mocleaner.exe <dir>
```

## License
See [LICENSE](LICENSE) file

## Development
 * Authors: [@dietmarw](https://github.com/dietmarw), [@mtiller](https://github.com/mtiller)
 * Contributors: See [graphs/contributors](../../graphs/contributors)

You may report any issues with using the [Issues](../../issues) button.

Contributions in shape of [Pull Requests](../../pulls) are always welcome.
