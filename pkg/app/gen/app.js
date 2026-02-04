// -----------------------------------------------------------------------------
// go-app
// -----------------------------------------------------------------------------
var goappNav = function () { };

var goappUpdatedBeforeWasmLoaded = false;
var goappOnUpdate = function () {
  goappUpdatedBeforeWasmLoaded = true;
};

var goappAppInstallChangedBeforeWasmLoaded = false;
var goappOnAppInstallChange = function () {
  goappAppInstallChangedBeforeWasmLoaded = true;
};

const goappEnv = {{.Env}};
const goappLoadingLabel = "{{.LoadingLabel}}";
const goappWasmContentLength = "{{.WasmContentLength}}";
const goappWasmContentLengthHeader = "{{.WasmContentLengthHeader}}";

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
      goappSetupPushNotification();
    } catch (err) {
      console.error("goapp service worker registration failed: ", err);
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
  registration.addEventListener("updatefound", (event) => {
    const newSW = registration.installing;
    newSW.addEventListener("statechange", (event) => {
      if (!navigator.serviceWorker.controller) {
        return;
      }

      switch (newSW.state) {
        case "activated":
          goappOnUpdate();
      }
    });
  });
}

function goappTryUpdate() {
  if (!goappServiceWorkerRegistration) {
    return;
  }
  goappServiceWorkerRegistration.update();
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
  return !goappIsAppInstalled() && (deferredPrompt != null || goappIsAppleBrowser());
}

function goappIsAppInstalled() {
  return navigator.standalone === true ||
    window.matchMedia("(display-mode: standalone)").matches ||
    document.referrer.startsWith('android-app://');
}

function goappIsAppleBrowser() {
  const ua = navigator.userAgent;
  const isIPadOS = /\bMacintosh\b/.test(ua) && navigator.maxTouchPoints > 1;
  const isIOSFamily = /iP(hone|ad|od)/.test(ua) || isIPadOS;
  const isMacSafari =
    /\bMacintosh\b/.test(ua) &&
    /\bSafari\b/.test(ua) &&
    !/\bChrome\b|\bEdg\b|\bOPR\b|\bBrave\b/.test(ua);
  return isIOSFamily || isMacSafari;
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

async function goappNewNotification(jsonNotification) {
  let notification = JSON.parse(jsonNotification);

  const title = notification.title;
  delete notification.title;

  let path = notification.path;
  if (!path) {
    path = "/";
  }

  if (!("serviceWorker" in navigator) || !goappServiceWorkerRegistration || !goappServiceWorkerRegistration.active) {
    const webNotification = new Notification(title, notification);

    webNotification.onclick = () => {
      goappNav(path);
      webNotification.close();
    };
    return;
  }

  const serviceWorker = goappServiceWorkerRegistration.active;
  serviceWorker.postMessage({
    type: "goapp:notify",
    options: notification,
  });
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
async function goappInitWebAssembly() {
  const loader = document.getElementById("app-wasm-loader");

  if (!goappCanLoadWebAssembly()) {
    loader.remove();
    return;
  }

  let instantiateStreaming = WebAssembly.instantiateStreaming;
  if (!instantiateStreaming) {
    instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const loaderIcon = document.getElementById("app-wasm-loader-icon");
  const loaderLabel = document.getElementById("app-wasm-loader-label");

  try {
    const showProgress = (progress) => {
      loaderLabel.innerText = goappLoadingLabel.replace("{progress}", progress);
    };
    showProgress(0);

    const go = new Go();
    const wasm = await instantiateStreaming(
      fetchWithProgress("{{.Wasm}}", showProgress),
      go.importObject
    );

    go.run(wasm.instance);
    loader.remove();
  } catch (err) {
    loaderIcon.className = "goapp-logo";
    loaderLabel.innerText = err;
    console.error("loading wasm failed: ", err);
  }
}

function goappCanLoadWebAssembly() {
  if (
    /bot|googlebot|crawler|spider|robot|crawling/i.test(navigator.userAgent)
  ) {
    return false;
  }

  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get("wasm") !== "false";
}

async function fetchWithProgress(url, progess) {
  const response = await fetch(url);

  let contentLength = goappWasmContentLength;
  if (contentLength <= 0) {
    try {
      contentLength = response.headers.get(goappWasmContentLengthHeader);
    } catch { }
    if (!goappWasmContentLengthHeader || !contentLength) {
      contentLength = response.headers.get("Content-Length");
    }
  }

  const total = parseInt(contentLength, 10);
  let loaded = 0;

  const progressHandler = function (loaded, total) {
    progess(Math.round((loaded * 100) / total));
  };

  var res = new Response(
    new ReadableStream(
      {
        async start(controller) {
          var reader = response.body.getReader();
          for (; ;) {
            var { done, value } = await reader.read();

            if (done) {
              progressHandler(total, total);
              break;
            }

            loaded += value.byteLength;
            progressHandler(loaded, total);
            controller.enqueue(value);
          }
          controller.close();
        },
      },
      {
        status: response.status,
        statusText: response.statusText,
      }
    )
  );

  for (var pair of response.headers.entries()) {
    res.headers.set(pair[0], pair[1]);
  }

  return res;
}
