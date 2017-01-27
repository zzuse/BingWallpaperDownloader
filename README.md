
# BingWallpaperDownloader   
written in golang.   
download bing beautiful wallpapers for only desktop usage.   
inspired by https://github.com/redstoneleo/BingNiceWallpapers   
inspired by https://blog.hyperexpert.com/shell-script-to-download-bings-daily-image-as-a-wallpaper-for-mac/   

Download 7 wallpapers from 7 market place one time
------------
If you run the program, it will download nearest 7 bing wallpapers from the Bing websites.   
Default it will save files to ./downloads directory.   
You can specify other directory to save.   
The filename is same as the Bing website, trim before the http slashes.   
Such as: VillersAbbey_ZH-CN10373383330_1920x1080.jpg   
Note: default resolution is 1920x1200, but U can specify parameters on your own. and saving image may fail cause some picture only have 1920x1080    

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
picture checksum history to avoid multiple download same file.   
more testing   
knowed issue: if U run this in China, maybe only get one zh-CN location, cause firewall   

#License
----------
WTFPL version 2



