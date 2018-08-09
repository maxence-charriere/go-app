
var goapp = {
    nodes: {}
};

function render(changes = []) {
    changes.forEach(c => {
        switch (c.Type) {
            case 'createText':
                createText(c);
                break;

            case 'setText':
                setText(c);
                break;

            case 'createElem':
                createElem(c);
                break;

            case 'setAttrs':
                setAttrs(c);
                break;

            case 'appendChild':
                appendChild(c);
                break;

            case 'removeChild':
                removeChild(c);
                break;

            case 'replaceChild':
                replaceChild(c);
                break;

            case 'createCompo':
                createCompo(c);
                break;

            case 'setCompoRoot':
                setCompoRoot(c);
                break;

            case 'deleteNode':
                deleteNode(c);
                break;

            default:
                console.log('unknown change: ' + c.Type);
        }
    });
}

function createText(change = {}) {
    const { ID } = change.Value;
    const n = document.createTextNode("");

    n.ID = ID;
    goapp.nodes[ID] = n;
}

function setText(change = {}) {
    const { ID, Text } = change.Value;

    const n = goapp.nodes[ID];
    if (!n) {
        return;
    }

    n.nodeValue = Text;
}

function createElem(change = {}) {
    const { ID, TagName } = change.Value;
    const n = document.createElement(TagName);

    n.ID = ID;
    goapp.nodes[ID] = n;
}

function setAttrs(change = {}) {
    const { ID, Attrs } = change.Value;

    const n = goapp.nodes[ID];
    if (!n) {
        return;
    }

    const nAttrs = n.attributes;
    const toDelete = [];

    for (var i = 0; i < nAttrs.length; i++) {
        const name = nAttrs[i].name;

        if (Attrs[name] === undefined) {
            toDelete.push(name);
        }
    }

    toDelete.forEach(name => {
        n.removeAttribute(name);
    });

    for (var name in Attrs) {
        const curVal = n.getAttribute(name);
        const newVal = Attrs[name];

        if (name === 'value') {
            n.value = newVal;
            continue;
        }

        if (curVal !== newVal) {
            n.setAttribute(name, newVal);
        }
    }
}

function appendChild(change = {}) {
    const { ParentID, ChildID } = change.Value;

    const parent = goapp.nodes[ParentID];
    if (!parent) {
        return;
    }

    const child = childRoot(goapp.nodes[ChildID]);
    if (!child) {
        return;
    }

    parent.appendChild(child);
}

function removeChild(change = {}) {
    const { ParentID, ChildID } = change.Value;

    const parent = goapp.nodes[ParentID];
    if (!parent) {
        return;
    }

    const child = childRoot(goapp.nodes[ChildID]);
    if (!child) {
        return;
    }

    parent.removeChild(child);
}

function replaceChild(change = {}) {
    const { ParentID, ChildID, OldID } = change.Value;

    const parent = goapp.nodes[ParentID];
    if (!parent) {
        return;
    }

    const newChild = childRoot(goapp.nodes[ChildID]);
    if (!newChild) {
        return;
    }


    const oldChild = childRoot(goapp.nodes[OldID]);
    if (!oldChild) {
        return;
    }

    parent.replaceChild(newChild, oldChild);
}

function createCompo(change = {}) {
    const { ID, Name } = change.Value;

    const compo = {
        ID,
        Name,
        IsCompo: true
    }

    goapp.nodes[ID] = compo;
}

function setCompoRoot(change = {}) {
    const { ID, RootID } = change.Value;
    const compo = goapp.nodes[ID];

    compo.RootID = RootID;
}

function deleteNode(change = {}) {
    const { ID } = change.Value;
    delete goapp.nodes[ID];
}

function childRoot(node) {
    if (!node || !node.IsCompo) {
        return node;
    }

    return goapp.nodes[node.RootID];
}

function mapObject(obj) {
    var map = {};

    for (var field in obj) {
        const name = field[0].toUpperCase() + field.slice(1);
        const value = obj[field];
        const type = typeof value;

        switch (type) {
            case 'object':
                break;

            case 'function':
                break;

            default:
                map[name] = value;
                break;
        }
    }
    return map;
}

function callCompoHandler(compoID, target, src, event) {
    switch (event.type) {
        case 'change':
            onchangeToGolang(compoID, target, src, event);
            break;

        case 'drag':
        case 'dragstart':
        case 'dragend':
        case 'dragexit':
            onDragStartToGolang(compoID, target, src, event);
            break;

        case 'dragenter':
        case 'dragleave':
        case 'dragover':
        case 'drop':
            ondropToGolang(compoID, target, src, event);
            break;

        default:
            eventToGolang(compoID, target, src, event);
            break;
    }
}

function onchangeToGolang(compoID, target, src, event) {
    golangRequest(JSON.stringify({
        'compo-id': compoID,
        'target': target,
        'json-value': JSON.stringify(src.value)
    }));
}

function onDragStartToGolang(compoID, target, src, event) {
    const payload = mapObject(event.dataTransfer);
    payload['Data'] = src.dataset.drag;

    event.dataTransfer.setData('text', src.dataset.drag);

    golangRequest(JSON.stringify({
        'compo-id': compoID,
        'target': target,
        'json-value': JSON.stringify(payload)
    }));
}

function ondropToGolang(compoID, target, src, event) {
    event.preventDefault();

    const payload = mapObject(event.dataTransfer);
    payload['Data'] = event.dataTransfer.getData('text');
    payload['file-override'] = 'xxx';

    golangRequest(JSON.stringify({
        'compo-id': compoID,
        'target': target,
        'json-value': JSON.stringify(payload),
        'override': 'Files'
    }));
}

function eventToGolang(compoID, target, src, event) {
    const payload = mapObject(event);

    if (src.contentEditable === 'true') {
        payload['InnerText'] = src.innerText;
    }

    golangRequest(JSON.stringify({
        'compo-id': compoID,
        'target': target,
        'json-value': JSON.stringify(payload)
    }));
}
