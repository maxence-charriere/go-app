const cacheName = "app-" + "{{.Version}}";

self.addEventListener("install", event => {
  console.log("installing app worker");
  self.skipWaiting();

  event.waitUntil(
    caches.open(cacheName).then(cache => {
      return cache.addAll([
        {{range $path, $element := .ResourcesToCache}}"{{$path}}",
        {{end}}
      ]);
    })
  );
});

self.addEventListener("activate", event => {
  console.log("app worker is activating");
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
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
