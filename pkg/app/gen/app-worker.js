const cacheName = "app-" + "{{.Version}}";

self.addEventListener("install", event => {
  console.log("installing app worker {{.Version}}");

  event.waitUntil(
    caches.open(cacheName).
      then(cache => {
        return cache.addAll([
          {{range $path, $element := .ResourcesToCache}}"{{$path}}",
          {{end}}
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
  console.log("app worker {{.Version}} is activated");
});

self.addEventListener("fetch", event => {
  event.respondWith(
    caches.match(event.request).then(response => {
      return response || fetch(event.request);
    })
  );
});
