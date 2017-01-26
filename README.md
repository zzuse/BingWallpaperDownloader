
# BingWallpaperDownloader   
written in golang.
download bing beautiful wallpapers for only desktop usage.
inspired by https://github.com/redstoneleo/BingNiceWallpapers

Download 7 wallpapers one time
------------
If you run the program, it will downloads nearest 7  bing wallpaper from the Bing websites(location:zh-cn).   
Default it will save files to ./downloads directory
You can specify other directory to save.
The filename is same as the Bing website, trim before the http slashes.
such as: VillersAbbey_ZH-CN10373383330_1920x1080.jpg

Build and Installation
------------
```shell
go build BingWallpaperDownload.go
```
for other platform usage binary cross compile: 
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
```
- - - - --
#TODO
other country location
filename add date
more testing

#License
----------
WTFPL version 2



