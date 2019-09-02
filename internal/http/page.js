// -----------------------------------------------------------------------------
// Init service worker
// -----------------------------------------------------------------------------
if ('serviceWorker' in navigator) {
  navigator.serviceWorker
    .register('/goapp.js')
    .then(reg => {
      console.log('offline service worker registered')
    })
    .catch(err => {
      console.error('offline service worker registration failed', err)
    })
}

// -----------------------------------------------------------------------------
// Init progressive app
// -----------------------------------------------------------------------------
let deferredPrompt

window.addEventListener('beforeinstallprompt', (e) => {
  e.preventDefault()
  deferredPrompt = e
  console.log('beforeinstallprompt')
})

// -----------------------------------------------------------------------------
// Init Web Assembly
// -----------------------------------------------------------------------------
if (!WebAssembly.instantiateStreaming) {
  WebAssembly.instantiateStreaming = async (resp, importObject) => {
    const source = await (await resp).arrayBuffer()
    return await WebAssembly.instantiate(source, importObject)
  }
}

const go = new Go()

WebAssembly
  .instantiateStreaming(fetch('/goapp.wasm'), go.importObject)
  .then((result) => {
    go.run(result.instance)
  })
  .catch(err => {
    const loadingIcon = document.getElementById('App_LoadingIcon')
    loadingIcon.className = ''

    const loadingLabel = document.getElementById('App_LoadingLabel')
    loadingLabel.innerText = err
    console.error('loading wasm failed: ' + err)
  })
