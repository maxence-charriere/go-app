const etag = '{{.ETag}}'

self.addEventListener('install', function (e) {
  console.log('intalling goapp worker', etag)
  self.skipWaiting()

  e.waitUntil(
    caches.open('goapp').then(function (cache) {
      return cache.addAll([
        {{range .Paths}}'{{.}}',
        {{end}}'/'
      ])
    })
  )
})

self.addEventListener('activate', event => {
  console.log('new version activated')
})

self.addEventListener('fetch', event => {
  event.respondWith(
    caches
      .match(event.request)
      .then(response => {
        if (response) {
          console.log('fetch from cache')
          return response
        }

        console.log('fetch from network')
        // event.request.headers.set('If-None-Match', '"' + etag + '"')
        if (event.request.headers) {
          console.log('headers found:', event.request.headers)
        }

        return fetch(event.request)
      })
  )
})
