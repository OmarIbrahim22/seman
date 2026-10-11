// Wires the WebAssembly build of seman (cmd/seman-web) to an xterm.js terminal.
(function () {
  "use strict";

  var STORAGE_KEY = "seman.semester";
  var COLS = 140;
  var ROWS = 34;
  var CELL = 0.6; // JetBrains Mono advance width, as a fraction of font size

  var MIN_FONT = 11;
  var MIN_ROWS = 30;

  var win = document.getElementById("window");
  var fullButton = document.getElementById("fullscreen");
  var lineRatio = 1.25; // row height / font size, measured once the terminal is open
  var screen = document.getElementById("screen");
  var loading = document.getElementById("loading");
  var toast = document.getElementById("toast");
  var term = null;
  var ready = false;
  var toastTimer = 0;

  function showToast(text) {
    toast.textContent = text;
    toast.classList.add("show");
    clearTimeout(toastTimer);
    toastTimer = setTimeout(function () { toast.classList.remove("show"); }, 4000);
  }

  // The app's footers need the full width, so the column count is fixed and
  // the font size follows the available width. Below the smallest readable
  // size the terminal scrolls sideways instead of shrinking further.
  // In fullscreen the font may grow much larger and the row count follows the
  // height, so the app fills the screen.
  function applyLayout() {
    var full = win.classList.contains("full");
    var width = screen.clientWidth - 20;
    var size = Math.max(MIN_FONT, Math.min(full ? 30 : 15, Math.floor(width / (COLS * CELL))));
    var rows = ROWS;
    if (full) {
      var height = screen.clientHeight - 20;
      rows = Math.floor(height / Math.ceil(size * lineRatio));
      if (rows < MIN_ROWS) {
        // Wide but short screen: shrink the text until enough rows fit.
        size = Math.max(MIN_FONT, Math.floor(height / MIN_ROWS / lineRatio));
        rows = Math.floor(height / Math.ceil(size * lineRatio));
      }
      rows = Math.max(20, Math.min(60, rows));
    }
    if (term.options.fontSize !== size) term.options.fontSize = size;
    if (term.rows !== rows) term.resize(COLS, rows);
    if (ready) window.semanResize(term.cols, term.rows);
  }

  function setFull(on) {
    win.classList.toggle("full", on);
    document.body.classList.toggle("locked", on);
    fullButton.textContent = on ? "Exit fullscreen" : "Fullscreen";
    applyLayout();
    term.focus();
  }

  function toggleFull() {
    if (win.classList.contains("full")) {
      if (document.fullscreenElement) document.exitFullscreen();
      setFull(false);
      return;
    }
    // The CSS class alone already fills the window, which covers browsers
    // without the Fullscreen API (iPhone Safari).
    setFull(true);
    if (!win.requestFullscreen) return;
    win.requestFullscreen().then(function () {
      // Esc cancels dialogs in seman; where the browser allows it, keep Esc
      // for the app (the browser then exits fullscreen on a long press).
      if (navigator.keyboard && navigator.keyboard.lock) {
        navigator.keyboard.lock(["Escape"]).catch(function () {});
      }
    }).catch(function () {});
  }

  function send(data) {
    if (ready) window.semanInput(data);
  }

  function start() {
    term = new Terminal({
      cols: COLS,
      rows: ROWS,
      fontSize: 13,
      fontFamily: '"JetBrains Mono", ui-monospace, Menlo, Consolas, "DejaVu Sans Mono", monospace',
      cursorBlink: false,
      scrollback: 0,
      theme: { background: "#11111b", foreground: "#cdd6f4", cursor: "#cba6f7" }
    });
    term.open(document.getElementById("terminal"));
    var grid = term.element.querySelector(".xterm-screen");
    if (grid && grid.clientHeight) lineRatio = grid.clientHeight / ROWS / 13;
    applyLayout();
    term.onData(send);
    window.addEventListener("resize", applyLayout);
    fullButton.addEventListener("click", toggleFull);
    document.addEventListener("fullscreenchange", function () {
      // Leaving browser fullscreen (Esc, F11) also restores the page layout.
      if (!document.fullscreenElement && win.classList.contains("full")) setFull(false);
      else applyLayout();
    });

    // Called from Go.
    window.semanWrite = function (text) { term.write(text); };
    window.semanReady = function () {
      ready = true;
      loading.classList.add("hidden");
      window.semanResize(term.cols, term.rows);
    };
    window.semanExited = function (error) {
      if (error) {
        ready = false;
        showToast("seman stopped: " + error);
      } else {
        showToast("You quit seman. In a terminal that returns you to the shell; here it just starts again.");
      }
    };

    var go = new Go();
    var load = WebAssembly.instantiateStreaming
      ? WebAssembly.instantiateStreaming(fetch("seman.wasm", { cache: "no-cache" }), go.importObject)
      : fetch("seman.wasm", { cache: "no-cache" }).then(function (r) { return r.arrayBuffer(); })
          .then(function (bytes) { return WebAssembly.instantiate(bytes, go.importObject); });
    load.then(function (result) { go.run(result.instance); })
      .catch(function (error) { loading.textContent = "Could not load seman: " + error.message; });
  }

  document.querySelectorAll("button[data-key]").forEach(function (button) {
    button.addEventListener("click", function () {
      send(button.dataset.key);
      if (term) term.focus();
    });
  });

  document.getElementById("reset").addEventListener("click", function () {
    try { localStorage.removeItem(STORAGE_KEY); } catch (e) { /* storage blocked */ }
    location.reload();
  });

  // xterm measures the font when it opens, so wait for the web font first.
  var fonts = document.fonts && document.fonts.load
    ? Promise.all([document.fonts.load('14px "JetBrains Mono"'), document.fonts.load('bold 14px "JetBrains Mono"')])
    : Promise.resolve();
  fonts.then(start, start);
})();
