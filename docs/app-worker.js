// -----------------------------------------------------------------------------
// PWA
// -----------------------------------------------------------------------------
const cacheName = "app-" + "8a09cae95cc54f7ede06c043c00fa408649b2edb";
const resourcesToCache = ["https://raw.githubusercontent.com/maxence-charriere/go-app/master/docs/web/icon.png","https://pagead2.googlesyndication.com/pagead/js/adsbygoogle.js?client=ca-pub-1013306768105236","/web/js/prism.js","/web/documents/what-is-go-app.md","/web/documents/updates.md","/web/documents/home.md","/web/documents/home-next.md","/web/css/prism.css","/web/css/docs.css","/web/app.wasm","/wasm_exec.js","/manifest.webmanifest","/app.js","/app.css","/"];

self.addEventListener("install", async (event) => {
  try {
    console.log("installing app worker 8a09cae95cc54f7ede06c043c00fa408649b2edb");
    await installWorker();
    await self.skipWaiting();
  } catch (error) {
    console.error("error during installation:", error);
  }
});

async function installWorker() {
  const cache = await caches.open(cacheName);
  await cache.addAll(resourcesToCache);
}

self.addEventListener("activate", async (event) => {
  try {
    await deletePreviousCaches(); // Await cache cleanup
    await self.clients.claim(); // Ensure the service worker takes control of the clients
    console.log("app worker 8a09cae95cc54f7ede06c043c00fa408649b2edb is activated");
  } catch (error) {
    console.error("error during activation:", error);
  }
});

async function deletePreviousCaches() {
  keys = await caches.keys();
  keys.forEach(async (key) => {
    if (key != cacheName) {
      try {
        console.log("deleting", key, "cache");
        await caches.delete(key);
      } catch (err) {
        console.error("deleting", key, "cache failed:", err);
      }
    }
  });
}

self.addEventListener("fetch", (event) => {
  event.respondWith(fetchWithCache(event.request));
});

async function fetchWithCache(request) {
  const cachedResponse = await caches.match(request);
  if (cachedResponse) {
    return cachedResponse;
  }
  return await fetch(request);
}

// -----------------------------------------------------------------------------
// Push Notifications
// -----------------------------------------------------------------------------
self.addEventListener("push", (event) => {
  if (!event.data || !event.data.text()) {
    return;
  }

  const notification = JSON.parse(event.data.text());
  if (!notification) {
    return;
  }

  const title = notification.title;
  delete notification.title;

  if (!notification.data) {
    notification.data = {};
  }
  let actions = [];
  for (let i in notification.actions) {
    const action = notification.actions[i];

    actions.push({
      action: action.action,
      path: action.path,
    });

    delete action.path;
  }
  notification.data.goapp = {
    path: notification.path,
    actions: actions,
  };
  delete notification.path;

  event.waitUntil(self.registration.showNotification(title, notification));
});

self.addEventListener("notificationclick", (event) => {
  event.notification.close();

  const notification = event.notification;
  let path = notification.data.goapp.path;

  for (let i in notification.data.goapp.actions) {
    const action = notification.data.goapp.actions[i];
    if (action.action === event.action) {
      path = action.path;
      break;
    }
  }

  event.waitUntil(
    clients
      .matchAll({
        type: "window",
      })
      .then((clientList) => {
        for (var i = 0; i < clientList.length; i++) {
          let client = clientList[i];
          if ("focus" in client) {
            client.focus();
            client.postMessage({
              goapp: {
                type: "notification",
                path: path,
              },
            });
            return;
          }
        }

        if (clients.openWindow) {
          return clients.openWindow(path);
        }
      })
  );
});
