// -----------------------------------------------------------------------------
// go-app
// -----------------------------------------------------------------------------
var goappNav = function () {};
var goappOnUpdate = function () {};
var goappOnAppInstallChange = function () {};

const goappEnv = {{.Env}};

let goappServiceWorkerRegistration;
let deferredPrompt = null;

goappInitServiceWorker();
goappWatchForUpdate();
goappWatchForInstallable();
goappInitWebAssembly();

// -----------------------------------------------------------------------------
// Service Worker
// -----------------------------------------------------------------------------
async function goappInitServiceWorker() {
  if ("serviceWorker" in navigator) {
    try {
      const registration = await navigator.serviceWorker.register(
        "{{.WorkerJS}}"
      );

      goappServiceWorkerRegistration = registration;
      goappSetupNotifyUpdate(registration);
      goappSetupAutoUpdate(registration);
      goappSetupPushNotification();
    } catch (err) {
      console.error("goapp service worker registration failed", err);
    }
  }
}

// -----------------------------------------------------------------------------
// Update
// -----------------------------------------------------------------------------
function goappWatchForUpdate() {
  window.addEventListener("beforeinstallprompt", (e) => {
    e.preventDefault();
    deferredPrompt = e;
    goappOnAppInstallChange();
  });
}

function goappSetupNotifyUpdate(registration) {
  registration.onupdatefound = () => {
    const installingWorker = registration.installing;

    installingWorker.onstatechange = () => {
      if (installingWorker.state != "installed") {
        return;
      }

      if (!navigator.serviceWorker.controller) {
        return;
      }

      goappOnUpdate();
    };
  };
}

function goappSetupAutoUpdate(registration) {
  const autoUpdateInterval = "{{.AutoUpdateInterval}}";
  if (autoUpdateInterval == 0) {
    return;
  }

  window.setInterval(() => {
    registration.update();
  }, autoUpdateInterval);
}

// -----------------------------------------------------------------------------
// Install
// -----------------------------------------------------------------------------
function goappWatchForInstallable() {
  window.addEventListener("appinstalled", () => {
    deferredPrompt = null;
    goappOnAppInstallChange();
  });
}

function goappIsAppInstallable() {
  return !goappIsAppInstalled() && deferredPrompt != null;
}

function goappIsAppInstalled() {
  const isStandalone = window.matchMedia("(display-mode: standalone)").matches;
  return isStandalone || navigator.standalone;
}

async function goappShowInstallPrompt() {
  deferredPrompt.prompt();
  await deferredPrompt.userChoice;
  deferredPrompt = null;
}

// -----------------------------------------------------------------------------
// Environment
// -----------------------------------------------------------------------------
function goappGetenv(k) {
  return goappEnv[k];
}

// -----------------------------------------------------------------------------
// Notifications
// -----------------------------------------------------------------------------
function goappSetupPushNotification() {
  navigator.serviceWorker.addEventListener("message", (event) => {
    const msg = event.data.goapp;
    if (!msg) {
      return;
    }

    if (msg.type !== "notification") {
      return;
    }

    goappNav(msg.path);
  });
}

async function goappSubscribePushNotifications(vapIDpublicKey) {
  try {
    const subscription =
      await goappServiceWorkerRegistration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: vapIDpublicKey,
      });
    return JSON.stringify(subscription);
  } catch (err) {
    console.error(err);
    return "";
  }
}

function goappNewNotification(jsonNotification) {
  let notification = JSON.parse(jsonNotification);

  const title = notification.title;
  delete notification.title;

  let path = notification.path;
  if (!path) {
    path = "/";
  }

  const webNotification = new Notification(title, notification);

  webNotification.onclick = () => {
    goappNav(path);
    webNotification.close();
  };
}

// -----------------------------------------------------------------------------
// Keep Clean Body
// -----------------------------------------------------------------------------
function goappKeepBodyClean() {
  const body = document.body;
  const bodyChildrenCount = body.children.length;

  const mutationObserver = new MutationObserver(function (mutationList) {
    mutationList.forEach((mutation) => {
      switch (mutation.type) {
        case "childList":
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
// Web Assembly
// -----------------------------------------------------------------------------

async function fetchWithProgress(path) {
  const response = await fetch(path);
  var contentLength = +response.headers.get('X-App-Wasm-Length');
  if (contentLength === 0) {
    contentLength = +response.headers.get('Content-Length');
  }
  const loaderLabel = document.getElementById("app-wasm-loader-label");
  let bytesLoaded = 0;
  const ts = new TransformStream({
    transform(chunk, ctrl) {
      bytesLoaded += chunk.byteLength;
      if (contentLength !== 0) {
        loaderLabel.innerText = `downloading ${(bytesLoaded / contentLength * 100).toFixed(2)}%`
      } else {
        loaderLabel.innerText = `downloading ${(bytesLoaded / 1000000).toFixed(2)} mb`
      }
      ctrl.enqueue(chunk)
    }
  });
  return new Response(response.body.pipeThrough(ts), response);
}

async function initWasmWithProgress(wasmFile, importObject) {
  if (typeof TransformStream === "function" && ReadableStream.prototype.pipeThrough) {
    let done = false;
    const response = await fetchWithProgress(wasmFile, function () {
      if (!done) {
        progress.apply(null, arguments);
      }
    });
    const wasm = await WebAssembly.instantiateStreaming(response, importObject);
    done = true;
    return wasm
  } else {
    // xhr fallback, this is slower and doesn't use WebAssembly.InstantiateStreaming,
    // but it's only happening on Firefox, and we can probably live with the app
    // starting slightly slower there...
    const loaderLabel = document.getElementById("app-wasm-loader-label");
    const xhr = new XMLHttpRequest();
    await new Promise(function (resolve, reject) {
      xhr.open("GET", wasmFile);
      xhr.responseType = "arraybuffer";
      xhr.onload = resolve;
      xhr.onerror = reject;
      xhr.onprogress = e => {
        if(e.lengthComputable) {
          loaderLabel.innerText = `downloading ${(e.loaded / e.total * 100).toFixed(2)}%`
        } else {
          loaderLabel.innerText = `downloading ${(e.loaded / 1000000).toFixed(2)} mb`
        }
      }
      xhr.send();
    });
    return await WebAssembly.instantiate(xhr.response, importObject);
  }
}

async function goappInitWebAssembly() {
  if (!goappCanLoadWebAssembly()) {
    document.getElementById("app-wasm-loader").style.display = "none";
    return;
  }

  try {
    const go = new Go();
    const wasm = await initWasmWithProgress("{{.Wasm}}", go.importObject);

    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "goapp-logo";

    go.run(wasm.instance);
  } catch (err) {
    const loaderIcon = document.getElementById("app-wasm-loader-icon");
    loaderIcon.className = "goapp-logo";

    const loaderLabel = document.getElementById("app-wasm-loader-label");
    loaderLabel.innerText = err;

    console.error("loading wasm failed: ", err);
  }
}

function goappCanLoadWebAssembly() {
  return !/bot|googlebot|crawler|spider|robot|crawling/i.test(
    navigator.userAgent
  );
}
