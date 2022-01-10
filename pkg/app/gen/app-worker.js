const cacheName = "app-" + "T(Version)";

self.addEventListener("install", event => {
  console.log("installing app worker " + "T(Version)");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll("T(ResourcesToCache)");
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
  console.log("app worker " + "T(Version)" + " is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
