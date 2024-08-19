<h1>
Downgram<br>
<small><small>a telegram media downloader</small></small>
</h1>

This is a small tool that helps download Telegram media files.

**Must notice: This is sufficient for normal use, but cannot guarantee 100% robustness.** I developed it as an auxiliary tool for use. When encountering an error, you can try restarting or trying again

**Requirements**
  - Theoretically supports mainstream platforms, but only tested on Windows 11 64bit

## Build

```
go mod tidy
go mod vendor

sh scripts/build.sh
```

Built files will be found in `./dist`

## Licenses

Downgram is licensed under the MIT License.

Downgram includes the following third-party libraries:

- gioui: Licensed under the MIT license. No changes has been made.
- gotd/td: Licensed under the MIT License. No changes has been made.
- ncruces/zenity: Licensed under the MIT license. No changes has been made.
