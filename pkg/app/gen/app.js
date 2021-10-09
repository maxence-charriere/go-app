// -----------------------------------------------------------------------------
// Init service worker
// -----------------------------------------------------------------------------
var goappOnUpdate = function () { };

if ("serviceWorker" in navigator) {
  navigator.serviceWorker
    .register("{{.WorkerJS}}")
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
const goappEnv = {{.Env }};

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
async function initWebAssembly(){
  if (!/bot|googlebot|crawler|spider|robot|crawling/i.test(navigator.userAgent)) {
    const go = new Go();
  
    let response = await fetch("{{.Wasm}}");
    const reader = response.body.getReader();
    const contentLength = +response.headers.get('Content-Length');
    let receivedLength = 0;
    let chunks = [];
    const loaderLabel = document.getElementById("app-wasm-loader-label");
    const loaderLabelText = loaderLabel.innerText;
    while (true) {
        const { done, value } = await reader.read();
        if (done) {
            break;
        }
        chunks.push(value);
        receivedLength += value.length;
        if (contentLength !== 0){
            loaderLabel.innerText = `${loaderLabelText} received ${(receivedLength/1000000).toFixed(2)}/${(contentLength/1000000).toFixed(2)}MB ${(receivedLength/contentLength*100).toFixed(2)}%`
        }else{
            loaderLabel.innerText = `${loaderLabelText} received ${(receivedLength/1000000).toFixed(2)}MB`
        }
    }
    let chunksAll = new Uint8Array(receivedLength); // (4.1)
    let position = 0;
    for (let chunk of chunks) {
        chunksAll.set(chunk, position); // (4.2)
        position += chunk.length;
    }
  
    WebAssembly.instantiate(chunksAll.buffer, go.importObject)
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
