# flutter_licenses

A simple script to get a list of all packages referenced by a pubspec.lock file

It parses the pubspec.lock yaml file and fetches the [pub.dev](https://pub.dev/) page to extract
the specified license there. If a license could not be found, it will print a warning message and count the license
as `ERROR` in the final report table.

## Usage

```
$ go get -u github.com/plan3t-one/flutter_licenses
$ flutter_licenses path/to/pubspec.lock
```

## Example Output

```
+------------+---------+
|  LICENSE   | # FOUND |
+------------+---------+
| ERROR      |       5 |
| Apache 2.0 |      14 |
| BSD        |     109 |
| MIT        |      25 |
+------------+---------+
|   TOTAL    |   153   |
+------------+---------+
```