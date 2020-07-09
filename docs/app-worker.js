const cacheName = "app-" + "1cd03c9bc88326425b75a99bcc5ac0a02da00561";

self.addEventListener("install", event => {
  console.log("installing app worker 1cd03c9bc88326425b75a99bcc5ac0a02da00561");
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
  console.log("app worker 1cd03c9bc88326425b75a99bcc5ac0a02da00561 is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
