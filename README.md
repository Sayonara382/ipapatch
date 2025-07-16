# ipapatch
a cli tool to patch IPAs and their plugins, fixing problems with the share sheet, widgets, VPNs, and more!

it uses [zxPluginsInject](https://github.com/asdfzxcvbn/zxPluginsInject) by default, which is a rewrite of choco's original patch.

due to being written in go and having **native load command injection**, it's fully cross-compatible with macOS, linux, and iOS! (windows untested, but it should work correctly from v2.1.0)

you can find the latest binaries on the [release page](https://github.com/asdfzxcvbn/ipapatch/releases/latest).

# requirements
~~ipapatch only has 1 external dependency: the `zip` command. ipapatch actually could be adjusted to not require it, but it would break some pretty obscure tools and iOS internals, and `zip` is practically available everywhere, so this shouldn't be an issue.~~

v2.0.0 has removed the dependency on `zip`! i figured it would be nice to just test whether or not this will REALLY noticably break things. for now, ipapatch is a fully portable, native binary!

v2.1.0 fixed an issue where the output ipa wasnt able to be extracted by some signing apps and other tools. you can also pass `--zip` to use the `zip` command, although this isn't really required anymore

# usage
```bash
$ ipapatch --help
usage: ipapatch [-h/--help] --input <path> [--output <path] [--dylib <path>] [--inplace] [--noconfirm] [--plugins-only] [--zip] [--version]

flags:
  --input path      the path to the ipa file to patch
  --output path     the path to the patched ipa file to create
  --dylib path      the path to the dylib to use instead of the embedded zxPluginsInject
  --inplace         takes priority over --output, use this to overwrite the input file
  --noconfirm       skip interactive confirmation when not using --inplace, overwriting a file that already exists, etc
  --plugins-only    only inject into plugin binaries (not the main executable)
  --zip             use the zip cli tool to remove files (shouldnt be needed anymore)

info:
  -h, --help        show usage and exit
  --version         show version and exit
```

# credits
big thanks to:

- Chocolate Fluffy for the original IPA patcher tweak
- blacktop for [ipsw](https://github.com/blacktop/ipsw) and [go-macho](https://github.com/blacktop/go-macho), making the native load command injection possible
