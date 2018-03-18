// render replaces the node with the given id by the given component.
function render (payload) {
  const {id, component} = payload

  const selector = '[data-goapp-id="' + id + '"]'
  const elem = document.querySelector(selector)

  if (!elem) {
    return
  }
  elem.outerHTML = component
}

// render replaces the attributes of the node with the given id by the given
// attributes.
function renderAttributes (payload) {
  const {id, attributes} = payload

  if (!attributes) {
    return
  }

  const selector = '[data-goapp-id="' + id + '"]'
  const elem = document.querySelector(selector)

  if (!elem) {
    return
  }

  if (!elem.hasAttributes()) {
    return
  }
  const elemAttrs = elem.attributes

  // Remove missing attributes.
  for (var i = 0; i < elemAttrs.length; i++) {
    const name = elemAttrs[i].name

    if (name === 'data-goapp-id') {
      continue
    }

    if (attributes[name] === undefined) {
      elem.removeAttribute(name)
    }
  }

  // Set attributes.
  for (var name in attributes) {
    const currentValue = elem.getAttribute(name)
    const newValue = attributes[name]

    if (name === 'value') {
      elem.value = newValue
      continue
    }

    if (currentValue !== newValue) {
      elem.setAttribute(name, newValue)
    }
  }
}

function mapObject (obj) {
  var map = {}

  for (var field in obj) {
    const name = field[0].toUpperCase() + field.slice(1)
    const value = obj[field]
    const type = typeof (value)

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

function callGoEventHandler (compoID, target, src, event) {
  var payload = null

  switch (event.type) {
    case 'change':
      onchangeToGolang(compoID, target, src, event)
      break

    case 'drag':
    case 'dragstart':
    case 'dragend':
    case 'dragexit':
      onDragStartToGolang(compoID, target, src, event)
      break

    case 'dragenter':
    case 'dragleave':
    case 'dragover':
    case 'drop':
      ondropToGolang(compoID, target, src, event)
      break

    default:
      eventToGolang(compoID, target, src, event)
      break
  }
}

function onchangeToGolang (compoID, target, src, event) {
  golangRequest(JSON.stringify({
    'CompoID': compoID,
    'Target': target,
    'JSONValue': JSON.stringify(src.value)
  }))
}

function onDragStartToGolang (compoID, target, src, event) {
  const payload = mapObject(event.dataTransfer)
  payload['Data'] = src.dataset.drag

  event.dataTransfer.setData('text', src.dataset.drag)

  golangRequest(JSON.stringify({
    'CompoID': compoID,
    'Target': target,
    'JSONValue': JSON.stringify(payload)
  }))
}

function ondropToGolang (compoID, target, src, event) {
  event.preventDefault()

  const payload = mapObject(event.dataTransfer)
  payload['Data'] = event.dataTransfer.getData('text')
  payload['FileOverride'] = 'xxx'

  golangRequest(JSON.stringify({
    'CompoID': compoID,
    'Target': target,
    'JSONValue': JSON.stringify(payload),
    'Override': 'Files'
  }))
}

function eventToGolang (compoID, target, src, event) {
  const payload = mapObject(event)

  if (src.contentEditable === 'true') {
    payload['InnerText'] = src.innerText
  }

  golangRequest(JSON.stringify({
    'CompoID': compoID,
    'Target': target,
    'JSONValue': JSON.stringify(payload)
  }))
}
