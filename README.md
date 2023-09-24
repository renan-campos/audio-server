# 🎵 audio-server

## v0 Demo:
```
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl -u joe:secret -d "" http://localhost:1323/v0/admin/audio
Audio resource created: "fdf6c9c2-75c5-4fac-91f0-3042958de58e"
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl -u joe:secret -d name="guitar" http://localhost:1323/v0/admin/audio/fdf6c9c2-75c5-4fac-91f0-3042958de58e
Uploaded metadata for "fdf6c9c2-75c5-4fac-91f0-3042958de58e":
{
        name: "guitar"
}
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl -u joe:secret -d "?" http://localhost:1323/v0/admin/audio/fdf6c9c2-75c5-4fac-91f0-3042958de58e/ogg
Uploaded ogg file for "fdf6c9c2-75c5-4fac-91f0-3042958de58e"
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl http://localhost:1323/v0/audio
List of audio metadata
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl http://localhost:1323/v0/audio/fdf6c9c2-75c5-4fac-91f0-3042958de58e
List of audio "fdf6c9c2-75c5-4fac-91f0-3042958de58e" metadata
[rcampos@rcampos-thinkpadt14sgen2i audio-server]$ curl http://localhost:1323/v0/audio/fdf6c9c2-75c5-4fac-91f0-3042958de58e/ogg
ogg file of "fdf6c9c2-75c5-4fac-91f0-3042958de58e"
```
