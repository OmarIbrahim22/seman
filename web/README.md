# Web demo

A landing page that runs the real seman in the browser. The Go code is compiled
to WebAssembly and drawn in an [xterm.js](https://xtermjs.org) terminal, so the
demo always behaves like the terminal app.

## Build and preview

```bash
./web/build.sh
python3 -m http.server -d web/site 8000
```

Then open <http://localhost:8000>. The page must be served over HTTP; opening
`index.html` from disk does not work because browsers will not load wasm from
`file://`.

## How it works

- `cmd/seman-web` is a second entry point, built only for `GOOS=js GOARCH=wasm`.
  It starts the same `app.Model` as `cmd/seman`, reads keys from xterm.js,
  writes frames back to it, and stores data in `localStorage` through the
  `storage.Store` interface.
- A first-time visitor gets a sample semester placed around today's date
  (`cmd/seman-web/demo.go`). "Reset demo data" on the page loads it again.
- Quitting restarts the app, since a browser tab has no shell to return to.
- The Lofi player needs `mpv` and `yt-dlp`, so it does not work in the browser.

## Why there is a `third_party` folder

Bubble Tea v0.25 and two of its dependencies do not compile for js/wasm. The
copies in `web/third_party` are the exact versions from `go.sum` with one added
`*_js.go` file each that stubs out the TTY, console and clipboard calls.
`web/go.mod` is the normal `go.mod` plus `replace` lines pointing at those
copies; the build passes it with `-modfile`, so the regular `go.mod`, and
`go install ...@latest`, are not affected. When dependencies change in
`go.mod`, copy the change into `web/go.mod` as well.

## Hosting

`web/site` is a plain static folder once built, so any static host works.
`seman.wasm` and `wasm_exec.js` are build output and are not committed; run
`./web/build.sh` before uploading.

## Licences

xterm.js (MIT) is vendored in `web/site/vendor`. The modules in
`web/third_party` keep their own licence files.
