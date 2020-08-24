const cacheName = "app-" + "f58036df4e5ce98e724407beb82b8e6525125e2f";

self.addEventListener("install", event => {
  console.log("installing app worker f58036df4e5ce98e724407beb82b8e6525125e2f");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        "/go-app",
        "/go-app/app.css",
        "/go-app/app.js",
        "/go-app/manifest.json",
        "/go-app/wasm_exec.js",
        "/go-app/web/app.wasm",
        "/go-app/web/docs.css",
        "https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
        "https://fonts.googleapis.com/css2?family=Roboto&display=swap",
        "https://storage.googleapis.com/murlok-github/icon-192.png",
        "https://storage.googleapis.com/murlok-github/icon-512.png",
        
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  event.waitUntil(
    caches.keys().then(keyList => {
      return Promise.all(
        keyList.map(key => {
          if (key !== cacheName) {
            return caches.delete(key);
          }
        })
      );
    })
  );
  console.log("app worker f58036df4e5ce98e724407beb82b8e6525125e2f is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
