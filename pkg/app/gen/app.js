// -----------------------------------------------------------------------------
// go-app Interop
// -----------------------------------------------------------------------------
var goappNav = function () {};
var goappOnUpdate = function () {};
var goappOnAppInstallChange = function () {};

const goappEnv = JSON.parse(`{{.Env}}`);

function goappGetenv(k) {
  return goappEnv[k];
}

// -----------------------------------------------------------------------------
// Service Worker
// -----------------------------------------------------------------------------
goappInitServiceWorker();

async function goappInitServiceWorker() {
  if ("serviceWorker" in navigator) {
    try {
      const registration = await navigator.serviceWorker.register(
        "{{.WorkerJS}}"
      );

      goappSetupUpdate(registration);
      goappSetupAutoUpdate(registration);
      goappRegisterSubscription(registration);
    } catch (err) {
      console.error("goapp service worker registration failed", err);
    }
  }
}

async function goappInitServiceWorker() {
  if ("serviceWorker" in navigator) {
    try {
      const registration = await navigator.serviceWorker.register(
        "{{.WorkerJS}}"
      );

      goappSetupUpdate(registration);
      goappSetupAutoUpdate(registration);
      goappRegisterSubscription(registration);
    } catch (err) {
      console.error("goapp service worker registration failed", err);
    }
  }
}

function goappSetupUpdate(registration) {
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

async function goappRegisterSubscription(registration) {
  const vapIDPublicKey = "{{.PushNotifications.VAPIDPublicKey}}";
  const registrationURL = "{{.PushNotifications.SubscriptionPayloadFormat}}";
  if (!vapIDPublicKey || !registrationURL) {
    return;
  }

  const subscription = await registration.pushManager.subscribe({
    userVisibleOnly: true,
    applicationServerKey: vapIDPublicKey,
  });

  console.log(subscription);

  let body = "{{.PushNotifications.SubscriptionPayloadFormat}}";
  body = body.replace("%s", JSON.stringify(subscription));

  fetch("{{.PushNotifications.RegistrationURL}}", {
    method: "post",
    headers: {
      "Content-type": "application/json",
    },
    body: body,
  });
}

// -----------------------------------------------------------------------------
// App install
// -----------------------------------------------------------------------------
let deferredPrompt = null;

window.addEventListener("beforeinstallprompt", (e) => {
  e.preventDefault();
  deferredPrompt = e;
  goappOnAppInstallChange();
});

window.addEventListener("appinstalled", () => {
  deferredPrompt = null;
  goappOnAppInstallChange();
});

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
// Keep body clean
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
// Init Web Assembly
// -----------------------------------------------------------------------------
if (!/bot|googlebot|crawler|spider|robot|crawling/i.test(navigator.userAgent)) {
  if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
      const source = await (await resp).arrayBuffer();
      return await WebAssembly.instantiate(source, importObject);
    };
  }

  const go = new Go();

  WebAssembly.instantiateStreaming(fetch("{{.Wasm}}"), go.importObject)
    .then((result) => {
      const loaderIcon = document.getElementById("app-wasm-loader-icon");
      loaderIcon.className = "goapp-logo";

      go.run(result.instance);
    })
    .catch((err) => {
      const loaderIcon = document.getElementById("app-wasm-loader-icon");
      loaderIcon.className = "goapp-logo";

      const loaderLabel = document.getElementById("app-wasm-loader-label");
      loaderLabel.innerText = err;

      console.error("loading wasm failed: " + err);
    });
} else {
  document.getElementById("app-wasm-loader").style.display = "none";
}

// -----------------------------------------------------------------------------
// Notifications
// -----------------------------------------------------------------------------
function goappNewNotification(notification) {
  goappShowNotification((title, options) => {
    try {
      const notification = new Notification(title, options);

      notification.onclick = (e) => {
        let target = options.target;
        if (!target) {
          target = "/";
        }

        goappNav(target);
        notification.close();
      };
    } catch (err) {
      console.log(err);
    }
  }, notification);
}

function goappShowNotification(showNotification, notification) {
  console.log(notification);

  const title = notification.title;
  delete notification.title;

  for (let action in notification.actions) {
    delete action.target;
  }

  showNotification(title, notification);
}
