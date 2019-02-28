const etag = '6dbf9374cd238f5692099816e582e488a265957e'

self.addEventListener('install', function (e) {
  console.log('intalling goapp worker', etag)
  self.skipWaiting()

  e.waitUntil(
    caches.open('goapp').then(function (cache) {
      return cache.addAll([
        '/goapp.wasm',
        '/hello.css',
        '/logo.png',
        '/space.jpg',
        '/wasm_exec.js',
        '/'
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
        return response || fetch(event.request)
      })
  )
})
