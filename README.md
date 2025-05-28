# IPAPatch
A CLI tool to patch IPAs and their plugins, fixing problems with the share sheet, widgets, VPNs, and more!

It uses [zxPluginsInject](https://github.com/asdfzxcvbn/zxPluginsInject) by default, which is a rewrite of choco's original patch.

Due to being written in Go and having **native load command injection**, it's fully cross-compatible with macOS, Linux, and **Windows**!

You can find the latest binaries on the [release page](https://github.com/Sayonara382/ipapatch/releases/latest).

# Requirements
- On **macOS** and **Linux**, ipapatch requires the `zip` command to be installed (available by default on most systems).
- On **Windows**, no external zip command is required—everything works out of the box.

# Usage
```bash
$ ipapatch --help
Usage: ipapatch [-h/--help] --input <path> [--output <path] [--dylib <path>] [--inplace] [--noconfirm] [--plugins-only] [--version]

Flags:
  --input path      The path to the ipa file to patch
  --output path     The path to the patched ipa file to create
  --dylib path      The path to the dylib to use instead of the embedded zxPluginsInject
  --inplace         Takes priority over --output, use this to overwrite the input file
  --noconfirm       Skip interactive confirmation when not using --inplace, overwriting a file that already exists, etc
  --plugins-only    Only inject into plugin binaries (not the main executable)

Info:
  -h, --help        Show usage and exit
  --version         Show version and exit
```

# Credits
Big thanks to:

- Chocolate Fluffy for the original IPA patcher tweak.
- blacktop for [ipsw](https://github.com/blacktop/ipsw) and [go-macho](https://github.com/blacktop/go-macho), making the native load command injection possible.
- asdfzxcvbn for [ipapatch](https://github.com/asdfzxcvbn/ipapatch).