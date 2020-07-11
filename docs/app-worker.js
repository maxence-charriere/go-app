const cacheName = "app-" + "6f47b93f81f7ccf25cf7f07756978562ccfdd831";

self.addEventListener("install", event => {
  console.log("installing app worker 6f47b93f81f7ccf25cf7f07756978562ccfdd831");
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
  console.log("app worker 6f47b93f81f7ccf25cf7f07756978562ccfdd831 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
