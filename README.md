# ares-deploy
Simple wrapper of [webOS ares commands](http://developer.lge.com/webOSTV/sdk/web-sdk/webos-tv-cli/using-webos-tv-cli/). It's inculde functions of:
- `ares-package`
- `ares-install`
- `ares-launch`

and also open browser at web inspector page.

## Dependences
Require [webOS TV CLI](http://developer.lge.com/webOSTV/sdk/web-sdk/webos-tv-cli/using-webos-tv-cli/) installed in `PATH`

## Usage
`ares-deploy --path="/path/to/source" --output="/path/where/to/output/ipk" --device="gopherTV" --browser="chrome" --list="false"`

or just 

`ares-deploy`

from your project folder(in case your device name is `webOs`)

## Flags
- `path`(shortcut `p`) Path with source code. Default value - current folder.
- `output`(shortcut `o`) Path where output ipk file. Default `path/build/` folder.
- `device`(shortcut `d`) Device to connect. Default `webOs`.
- `browser`(shortcut `b`) Brower where open web inspector. Default browser of your sistem.
- `list`(shortcut `l`) Print available devices in json format. Default `false`.