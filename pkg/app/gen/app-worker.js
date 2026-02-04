// -----------------------------------------------------------------------------
// PWA
// -----------------------------------------------------------------------------
const cacheName = "app-" + "{{.Version}}";
const resourcesToCache = {{.ResourcesToCache}};

self.addEventListener("install", async (event) => {
  try {
    console.log("installing app worker {{.Version}}");
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
    console.log("app worker {{.Version}} is activated");
  } catch (error) {
    console.error("error during activation:", error);
  }
});

async function deletePreviousCaches() {
  const keys = await caches.keys();
  await Promise.all(
    keys.map(async (key) => {
      if (key !== cacheName) {
        try {
          console.log("deleting", key, "cache");
          await caches.delete(key);
        } catch (err) {
          console.error("deleting", key, "cache failed:", err);
        }
      }
    })
  );
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

  event.waitUntil(
    showNotification(self.registration, notification)
  );
});

self.addEventListener("message", (event) => {
  const msg = event.data;
  if (!msg || msg.type !== "goapp:notify") {
    return;
  }

  event.waitUntil(
    showNotification(self.registration, msg.options)
  );
});

async function showNotification(registration, notification) {
  const title = notification.title || "";
  // let delay = notification.delay || 0;

  let actions = [];
  for (let i in notification.actions) {
    const action = notification.actions[i];
    actions.push({
      action: action.action,
      path: action.path,
    });
    delete action.path;
  }

  notification.data = notification.data || {};
  notification.data.goapp = {
    path: notification.path,
    actions: actions,
  };
  delete notification.title;
  delete notification.path;
  delete notification.delay;

  // if (delay > 0) {
  //   delay = Math.floor(delay / 1e6);
  //   await sleep(delay);
  // }
  await registration.showNotification(title, notification);
}

function sleep(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

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