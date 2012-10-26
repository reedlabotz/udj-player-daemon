udj-player-daemon
=================

A daemon to run on peoples computers allowing the udj web client to actually play 
the music and do all communication with the server.

API
---
`/list?path=/some/path` 
Returns a json list of all files and directories at the given path. Each entry has the following format
```json
{
  "Name": "file or directory name",
  "IsDir": true
}
```

`/music?path=/path/to/music/folder` 
Returns a json list of all music files found recursivly in give path. Each item has the following format
```json
{
  "Title": "That song",
  "Artist": "Artist",
  "Album": "Album",
  "Length": "203",
  "Path": "/path/to/file.mp3"
}
```

`/stream?path=/path/to/music/folder/file.mp3` 
Returns the raw mp3 data to be played by the browser