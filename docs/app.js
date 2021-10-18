// -----------------------------------------------------------------------------
// Init service worker
// -----------------------------------------------------------------------------
var goappOnUpdate = function () { };

if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("/app-worker.js")
    .then(reg => {
      console.log("registering app service worker");

      reg.onupdatefound = function () {
        const installingWorker = reg.installing;
        installingWorker.onstatechange = function () {
          if (installingWorker.state == "installed") {
            if (navigator.serviceWorker.controller) {
              goappOnUpdate();
            }
          }
        };
      }
    })
    .catch(err => {
      console.error("offline service worker registration failed", err);
    });
}

// -----------------------------------------------------------------------------
// Env
// -----------------------------------------------------------------------------
const goappEnv = {"GOAPP_INTERNAL_URLS":"null","GOAPP_ROOT_PREFIX":"","GOAPP_STATIC_RESOURCES_URL":"","GOAPP_VERSION":"8e0f1317d923632879808cfe869d42560fd3eecd"};

function goappGetenv(k) {
  return goappEnv[k];
}

// -----------------------------------------------------------------------------
// App install
// -----------------------------------------------------------------------------
let deferredPrompt = null;
var goappOnAppInstallChange = function () { };

window.addEventListener("beforeinstallprompt", e => {
  e.preventDefault();
  deferredPrompt = e;
  goappOnAppInstallChange();
});

window.addEventListener('appinstalled', () => {
  deferredPrompt = null;
  goappOnAppInstallChange();
});

function goappIsAppInstallable() {
  return !goappIsAppInstalled() && deferredPrompt != null;
}

function goappIsAppInstalled() {
  const isStandalone = window.matchMedia('(display-mode: standalone)').matches;
  return isStandalone || navigator.standalone;
}

async function goappShowInstallPrompt() {
  deferredPrompt.prompt();
  await deferredPrompt.userChoice;
  deferredPrompt = null;
}

// -----------------------------------------------------------------------------
// Keep body clean
// -----------------------------------------------------------------------------
function goappKeepBodyClean() {
  const body = document.body;
  const bodyChildrenCount = body.children.length;

  const mutationObserver = new MutationObserver(function (mutationList) {
    mutationList.forEach((mutation) => {
      switch (mutation.type) {
        case 'childList':
          while (body.children.length > bodyChildrenCount) {
            body.removeChild(body.lastChild);
          }
          break;
      }
    });
  });

  mutationObserver.observe(document.body, {
    childList: true,
  });

  return () => mutationObserver.disconnect();
}

// -----------------------------------------------------------------------------
// Init Web Assembly
// -----------------------------------------------------------------------------
async function initWebAssembly() {
  if (!/bot|googlebot|crawler|spider|robot|crawling/i.test(navigator.userAgent)) {
    const go = new Go();

    const loaderLabel = document.getElementById("app-wasm-loader-label");

    let response = await fetch("/web/app.wasm");
    const reader = response.body.getReader();
    const contentLength = response.headers.get('Content-Length');

    let chunks = [];
    let len = 0;

    while (true) {
      const { done, value } = await reader.read();
      if (done) {
        break;
      }

      chunks.push(value);
      len += value.length;
      let progress = (len * 100 / contentLength).toFixed(2)
      if (progress > 100) {
        progress = (100).toFixed(2);
      }
      loaderLabel.innerText = progress + "%";
    }

    let wasmFile = new Uint8Array(len);
    let idx = 0;
    for (let c of chunks) {
      wasmFile.set(c, idx);
      idx += c.length;
    }

    WebAssembly.instantiate(wasmFile.buffer, go.importObject)
      .then(result => {
        const loaderIcon = document.getElementById("app-wasm-loader-icon");
        loaderIcon.className = "goapp-logo";

        go.run(result.instance);
      })
      .catch(err => {
        const loaderIcon = document.getElementById("app-wasm-loader-icon");
        loaderIcon.className = "goapp-logo";

        const loaderLabel = document.getElementById("app-wasm-loader-label");
        loaderLabel.innerText = err;

        console.error("loading wasm failed: " + err);
      });
  } else {
    document.getElementById('app-wasm-loader').style.display = "none";
  }
}

initWebAssembly();
