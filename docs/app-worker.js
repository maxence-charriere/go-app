const cacheName = "app-" + "8f9292a4cbed0e83e98b9154be3ae0ab8a3447cb";

self.addEventListener("install", event => {
  console.log("installing app worker 8f9292a4cbed0e83e98b9154be3ae0ab8a3447cb");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          "/",
          "/app.css",
          "/app.js",
          "/manifest.webmanifest",
          "/wasm_exec.js",
          "/web/app.wasm",
          "/web/css/docs.css",
          "/web/css/prism.css",
          "/web/documents/home-next.md",
          "/web/documents/home.md",
          "/web/documents/updates.md",
          "/web/documents/what-is-go-app.md",
          "/web/js/prism.js",
          "https://fonts.googleapis.com/css2?family=Montserrat:wght@400;500&display=swap",
          "https://storage.googleapis.com/murlok-github/icon-192.png",
          "https://storage.googleapis.com/murlok-github/icon-512.png",
          
        ]);
      }).
      then(() => {
        self.skipWaiting();
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
  console.log("app worker 8f9292a4cbed0e83e98b9154be3ae0ab8a3447cb is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
