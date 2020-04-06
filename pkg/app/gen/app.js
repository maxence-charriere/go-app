// -----------------------------------------------------------------------------
// Init service worker
// -----------------------------------------------------------------------------
if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("/app-worker.js")
    .then(reg => {
      console.log("registering app service worker");
    })
    .catch(err => {
      console.error("offline service worker registration failed", err);
    });
}

// -----------------------------------------------------------------------------
// Init progressive app
// -----------------------------------------------------------------------------
let deferredPrompt;
const goappEnv = {{.Env }}

window.addEventListener("beforeinstallprompt", e => {
  e.preventDefault();
  deferredPrompt = e;
});

// -----------------------------------------------------------------------------
// Init Web Assembly
// -----------------------------------------------------------------------------
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer();
    return await WebAssembly.instantiate(source, importObject);
  };
}

const go = new Go();

WebAssembly.instantiateStreaming(fetch("{{.Wasm}}"), go.importObject)
  .then(result => {
    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "app-wasm-icon";

    go.run(result.instance);
  })
  .catch(err => {
    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "app-wasm-icon";

    const loaderLabel = document.getElementById("app-wasm-loader-label");
    loaderLabel.innerText = err;

    console.error("loading wasm failed: " + err);
  });
