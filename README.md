
# BingWallpaperDownloader   
written in golang.   
download bing beautiful wallpapers for only desktop usage.   
inspired by https://github.com/redstoneleo/BingNiceWallpapers   
inspired by https://blog.hyperexpert.com/shell-script-to-download-bings-daily-image-as-a-wallpaper-for-mac/   

Download 1 wallpapers from 7 market places one time
------------
If you run the program and specify no parameters, it will download latest date bing wallpapers from 7 markets Bing websites.   
Default it will save files to ./downloads directory.   
The filename is same as the Bing website, trim before the http slashes.   
Such as: VillersAbbey_ZH-CN10373383330_1920x1080.jpg   
Default resolution is 1920x1200, but U can specify parameters on your own. and saving image may fail cause some picture only have 1920x1080    
You can specify -s to save to other directory.    
You can specify -d to download n days of wallpapers.   
You can specify -r to set resolution.   

Build and Installation
------------
default is on OSX:
```shell
go build BingWallpaperDownload.go
```
for other platforms usage binary cross compile: 
```shell
GOOS=linux GOARCH=amd64 go build -o BingWallpaperDownload BingWallpaperDownload.go  
GOOS=windows GOARCH=amd64 go build -o BingWallpaperDownload.exe BingWallpaperDownload.go
```

- - - - --
Run
------------
```shell
just run:
BingWallpaperDownload
or
BingWallpaperDownload.exe
more parameters specifications not documented yet.   
But you can have a try.   
```
- - - - --
#TODO
know issue: if the picture doesn't have 1920x1200, it will download a zero size picture. will fix this bug with the URL in json
picture checksum to avoid multiple download same file.   
more testing   
knowed issue: if U run this in China, maybe get only one zh-CN location, cause firewall   

#License
----------
WTFPL version 2



