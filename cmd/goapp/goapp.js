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
  console.log('goapp worker', etag, 'is activated')
})

self.addEventListener('fetch', event => {
  event.respondWith(
    caches
      .match(event.request)
      .then(response => {
        return response || fetch(event.request)
      })
  )
})
