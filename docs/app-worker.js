const cacheName = "app-" + "bb32bf73ca55e0eeb8f73666fe55dcf2ca8a3519";

self.addEventListener("install", event => {
  console.log("installing app worker bb32bf73ca55e0eeb8f73666fe55dcf2ca8a3519");
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
  console.log("app worker bb32bf73ca55e0eeb8f73666fe55dcf2ca8a3519 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
