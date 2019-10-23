const etag = '{{.ETag}}'

const goappCache = {
  name: 'goapp-cache',
  urls: [
    {{range. Paths}}'{{.}}',
    {{end}}'/'
  ]
}

self.addEventListener('install', function (event) {
  console.log('intalling goapp worker', etag)
  self.skipWaiting()

  event.waitUntil(
    caches.open(goappCache.name)
      .then(function (cache) {
        return cache.addAll(goappCache.urls)
      })
  )
})

self.addEventListener('fetch', function (event) {
  event.respondWith(
    caches.match(event.request)
      .then(function (response) {
        if (response) {
          return response
        }
        return fetch(event.request)
      })
  )
})

self.addEventListener('activate', function (event) {
  console.log('goapp worker', etag, 'is activated')

  const cacheWhitelist = [goappCache.name]

  event.waitUntil(
    caches.keys()
      .then(function (cacheNames) {
        return Promise.all(cacheNames.map(function (cacheName) {
          if (cacheWhitelist.indexOf(cacheNames) === -1) {
            return caches.delete(cacheNames)
          }
        }))
      })
  )
})
