function rand() {
  return Math.floor(Math.random() * 42 + 1);
}

function testRender() {
  render({
    id: 'test-01',
    component: `<p data-goapp-id="test-01">component ` + rand() + `</p>`
  });
}

function testRenderAttributes() {
  renderAttributes({
    id: 'test-02',
    attributes: {
      'data-new': 'added attribute',
      'data-update': rand()
    }
  });

  let component = `<ul data-goapp-id="test-02-bis">`;

  const elem = document.querySelector('[data-goapp-id="test-02"]');
  if (!elem.hasAttributes()) {
    return;
  }
  const elemAttrs = elem.attributes;

  for (let i = 0; i < elemAttrs.length; i++) {
    const { name, value } = elemAttrs[i];

    if (name === 'data-goapp-id') {
      continue;
    }

    component += '<li>' + name + ': ' + value + '</li>';
  }
  component += '</ul>';

  render({
    id: 'test-02-bis',
    component: component
  });
}

window.onload = () => {
  testMount();
};
