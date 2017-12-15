function testMount () {
  mount(`
<h1>Starting test</h1>

<h2>render</h2>
<button onclick="testRender()">Launch</button>
<h3>Output:</h3>
<p data-goapp-id="test-01"></p>

<h2>renderAttibutes</h2>
<button onclick="testRenderAttributes()">Launch</button>
<h3>Output:</h3>
<p data-goapp-id="test-02" data-remove="true" data-update="">
  <ul data-goapp-id="test-02-bis">
    <li>data-remove: true</li>
    <li>data-update:</li>
  </ul>
</p>

<h2>callback onclick</h2>
<button onclick="callGoEventHandler('compo-03', 'test', this, event)">Launch</button>
  `)
}

function rand () {
  return Math.floor((Math.random() * 42) + 1)
}

function testRender () {
  const component = `<p data-goapp-id="test-01">component ` + rand() + `</p>`
  render('test-01', component)
}

function testRenderAttributes () {
  renderAttributes('test-02', {
    'data-new': 'added attribute',
    'data-update': rand()
  })

  let component = `<ul data-goapp-id="test-02-bis">`

  const elem = document.querySelector('[data-goapp-id="test-02"]')
  if (!elem.hasAttributes()) {
    return
  }
  const elemAttrs = elem.attributes

  for (let i = 0; i < elemAttrs.length; i++) {
    const {name, value} = elemAttrs[i]

    if (name === 'data-goapp-id') {
      continue
    }

    component += '<li>' + name + ': ' + value + '</li>'
  }
  component += '</ul>'

  render('test-02-bis', component)
}

window.onload = () => {
  testMount()
}
