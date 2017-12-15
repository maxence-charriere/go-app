// mount mounts the given component in the page root.
function mount (component) {
  const selector = '[data-goapp-root]'
  const root = document.querySelector(selector)
  root.innerHTML = component
}

// render replaces the node with the given id by the given component.
function render (id, component) {
  const selector = '[data-goapp-id="' + id + '"]'
  const elem = document.querySelector(selector)
  elem.outerHTML = component
}

// render replaces the attributes of the node with the given id by the given
// attributes.
function renderAttributes (id, attrs) {
  if (!attrs) {
    return
  }

  const selector = '[data-goapp-id="' + id + '"]'
  const elem = document.querySelector(selector)

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

    if (attrs[name] === undefined) {
      elem.removeAttribute(name)
    }
  }

  // Set attributes.
  for (var name in attrs) {
    const currentValue = elem.getAttribute(name)
    const newValue = attrs[name]

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
