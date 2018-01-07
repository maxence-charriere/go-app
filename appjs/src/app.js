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

    if (currentValue !== newValue) {
      elem.setAttribute(name, newValue)
    }
  }
}

function callGoEventHandler (compoID, target, self, event) {
  var payload = {
    Value: self.value
  }

  for (var field in event) {
    const name = field[0].toUpperCase() + field.slice(1)
    const value = event[field]
    const type = typeof (value)

    switch (type) {
      case 'object':
        break

      case 'function':
        payload[name] = value.name
        break

      default:
        payload[name] = value
        break
    }
  }

  const mapping = {
    CompoID: compoID,
    Target: target,
    JSONValue: JSON.stringify(payload)
  }

  golangRequest(JSON.stringify(mapping))
}
