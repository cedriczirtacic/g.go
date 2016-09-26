# g

This is an stolen idea from https://github.com/rgbkrk/i (props to @rgbkrk)
There're some small changes:

1. Reads the bookmarks from an exported HTML file.
2. You can modify and reload the bookmarks file.
3. Using the "?print" query you can print all bookmarks to screen.
3. Runs with parameters like port number, disable commands (from query) or load another bookmark file (default is *g.html*).


## Build
```bash
$ go get "golang.org/x/net/html"
$ go build -o g g.go
```
If you get the package "golang.org/x/net/html" (pkg/ and src/) inside the same directory as g.go, make sure you set the **$GOPATH** environment variable to **$PWD** before the build.

## Run!
It's easy!
```bash
$ ./g
2016/09/25 21:31:55 Reading bookmarks from: g.html
2016/09/25 21:31:55 Listening on port: 8080
```
Or load another bookmarks file with port 80 (requires root perms) and disable commands.
```bash
$ sudo ./g -port 80 -disable_cmds -ifile test.html
2016/09/25 21:33:06 Reading bookmarks from: test.html
2016/09/25 21:33:06 Listening on port: 80
```
Do you have any firewall up? Check that before running this! :feelsgood:

## Commands
Now there are two commands: **print** and **reload**. Commands are passed thru the URL as a query string.

* "http://localhost:8080/**?print**": See what bookmarks you have.
```
Bookmarks:
	(*) cedriczirtacic!
```
* "http://localhost:8080/**?reload**": In case of a modification, reload all bookmarks.
```
Bookmarks reloaded!
```
