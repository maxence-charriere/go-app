// -----------------------------------------------------------------------------
// Goapp
// -----------------------------------------------------------------------------
var goapp = {
  nodes: {},

  actions: Object.freeze({
    'setRoot': 0,
    'newNode': 1,
    'delNode': 2,
    'setAttr': 3,
    'delAttr': 4,
    'setText': 5,
    'appendChild': 6,
    'removeChild': 7,
    'replaceChild': 8
  }),

  pointer: {
    x: 0,
    y: 0
  }
}

function render (changes = []) {
  changes.forEach(c => {
    switch (c.Action) {
      case goapp.actions.setRoot:
        setRoot(c)
        break

      case goapp.actions.newNode:
        newNode(c)
        break

      case goapp.actions.delNode:
        delNode(c)
        break

      case goapp.actions.setAttr:
        setAttr(c)
        break

      case goapp.actions.delAttr:
        delAttr(c)
        break

      case goapp.actions.setText:
        setText(c)
        break

      case goapp.actions.appendChild:
        appendChild(c)
        break

      case goapp.actions.removeChild:
        removeChild(c)
        break

      case goapp.actions.replaceChild:
        replaceChild(c)
        break

      default:
        console.log(c.Type + ' change is not supported')
    }
  })
}

function setRoot (change = {}) {
  const { NodeID } = change

  const n = goapp.nodes[NodeID]
  n.IsRootCompo = true

  const root = compoRoot(n)
  if (!root) {
    return
  }

  while (document.body.firstChild) {
    document.body.removeChild(document.body.firstChild)
  }

  document.body.appendChild(root)
}

function newNode (change = {}) {
  const { IsCompo = false, Type, NodeID, CompoID, Namespace } = change

  if (IsCompo) {
    goapp.nodes[NodeID] = {
      Type,
      ID: NodeID,
      IsCompo
    }

    return
  }

  var n = null

  if (Type === 'text') {
    n = document.createTextNode('')
  } else if (change.Namespace) {
    n = document.createElementNS(Namespace, Type)
  } else {
    n = document.createElement(Type)
  }

  n.ID = NodeID
  n.CompoID = CompoID
  goapp.nodes[NodeID] = n
}

function delNode (change = {}) {
  const { NodeID } = change
  delete goapp.nodes[NodeID]
}

function setAttr (change = {}) {
  const { NodeID, Key, Value = '' } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  n.setAttribute(Key, Value)
}

function delAttr (change = {}) {
  const { NodeID, Key } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  n.removeAttribute(Key)
}

function setText (change = {}) {
  const { NodeID, Value } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  n.nodeValue = Value
}

function appendChild (change = {}) {
  const { NodeID, ChildID } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  if (n.IsCompo) {
    n.RootID = ChildID
    return
  }

  const c = compoRoot(goapp.nodes[ChildID])
  if (!c) {
    return
  }

  n.appendChild(c)
}

function removeChild (change = {}) {
  const { NodeID, ChildID } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  const c = compoRoot(goapp.nodes[ChildID])
  if (!c) {
    return
  }

  n.removeChild(c)
}

function replaceChild (change = {}) {
  const { NodeID, ChildID, NewChildID } = change

  const n = goapp.nodes[NodeID]
  if (!n) {
    return
  }

  const c = compoRoot(goapp.nodes[ChildID])
  if (!c) {
    return
  }

  const nc = compoRoot(goapp.nodes[NewChildID])
  if (!nc) {
    return
  }

  if (n.IsCompo) {
    n.RootID = NewChildID

    if (n.IsRootCompo) {
      setRoot({ NodeID: n.ID })
    }

    return
  }

  n.replaceChild(nc, c)
}

function compoRoot (node) {
  if (!node || !node.IsCompo) {
    return node
  }

  const n = goapp.nodes[node.RootID]
  return compoRoot(n)
}

function mapObject (obj) {
  var map = {}

  for (var field in obj) {
    const name = field[0].toUpperCase() + field.slice(1)
    const value = obj[field]
    const type = typeof value

    switch (type) {
      case 'object':
        break

      case 'function':
        break

      default:
        map[name] = value
        break
    }
  }

  return map
}

function callCompoHandler (elem, event, fieldOrMethod) {
  switch (event.type) {
    case 'change':
      onchangeToGolang(elem, fieldOrMethod)
      break

    case 'drag':
    case 'dragstart':
    case 'dragend':
    case 'dragexit':
      onDragStartToGolang(elem, event, fieldOrMethod)
      break

    case 'dragenter':
    case 'dragleave':
    case 'dragover':
    case 'drop':
      ondropToGolang(elem, event, fieldOrMethod)
      break

    case 'contextmenu':
      event.preventDefault()
      trackPointerPosition(event)
      eventToGolang(elem, event, fieldOrMethod)
      break

    default:
      trackPointerPosition(event)
      eventToGolang(elem, event, fieldOrMethod)
  }
}

function onchangeToGolang (elem, fieldOrMethod) {
  goapp.emit(JSON.stringify({
    'CompoID': elem.CompoID,
    'FieldOrMethod': fieldOrMethod,
    'JSONValue': JSON.stringify(elem.value)
  }))
}

function onDragStartToGolang (elem, event, fieldOrMethod) {
  const payload = mapObject(event.dataTransfer)
  payload['Data'] = elem.dataset.drag
  setPayloadSource(payload, elem)

  event.dataTransfer.setData('text', elem.dataset.drag)

  goapp.emit(JSON.stringify({
    'CompoID': elem.CompoID,
    'FieldOrMethod': fieldOrMethod,
    'JSONValue': JSON.stringify(payload)
  }))
}

function ondropToGolang (elem, event, fieldOrMethod) {
  event.preventDefault()

  const payload = mapObject(event.dataTransfer)
  payload['Data'] = event.dataTransfer.getData('text')
  payload['FileOverride'] = 'xxx'
  setPayloadSource(payload, elem)

  goapp.emit(JSON.stringify({
    'CompoID': elem.CompoID,
    'FieldOrMethod': fieldOrMethod,
    'JSONValue': JSON.stringify(payload),
    'Override': 'Files'
  }))
}

function eventToGolang (elem, event, fieldOrMethod) {
  const payload = mapObject(event)
  setPayloadSource(payload, elem)

  if (elem.contentEditable === 'true') {
    payload['InnerText'] = elem.innerText
  }

  goapp.emit(JSON.stringify({
    'CompoID': elem.CompoID,
    'FieldOrMethod': fieldOrMethod,
    'JSONValue': JSON.stringify(payload)
  }))
}

function setPayloadSource (payload, elem) {
  payload['Source'] = {
    'GoappID': elem.ID,
    'CompoID': elem.CompoID,
    'ID': elem.id,
    'Class': elem.className,
    'Data': elem.dataset,
    'Value': elem.value
  }
}

function trackPointerPosition (event) {
  if (event.clientX != undefined) {
    goapp.pointer.x = event.clientX
  }

  if (event.clientY != undefined) {
    goapp.pointer.y = event.clientY
  }
}

// -----------------------------------------------------------------------------
// Context menu
// -----------------------------------------------------------------------------

function showContextMenu () {
  const bg = document.getElementById('App_ContextMenuBackground')
  if (!bg) {
    console.log('no context menu declared')
    return
  }
  bg.style.display = 'block'

  const menu = document.getElementById('App_ContextMenu')

  const width = window.innerWidth ||
    document.documentElement.clientWidth ||
    document.body.clientWidth

  const height = window.innerHeight ||
    document.documentElement.clientHeight ||
    document.body.clientHeight

  var x = goapp.pointer.x
  if (x + menu.offsetWidth > width) {
    x = width - menu.offsetWidth - 1
  }

  var y = goapp.pointer.y
  if (y + menu.offsetHeight > height) {
    y = height - menu.offsetHeight - 1
  }

  menu.style.left = x + 'px'
  menu.style.top = y + 'px'
}

function hideContextMenu () {
  const bg = document.getElementById('App_ContextMenuBackground')
  if (!bg) {
    console.log('no context menu declared')
    return
  }
  bg.style.display = 'none'
}

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
    console.error('wasm run failed: ' + err)
  })
