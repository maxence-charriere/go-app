
var goapp = {
    nodes: {
        "root:": document.body
    },

    actions: Object.freeze({
        "setRoot": 0,
        "newNode": 1,
        "delNode": 2,
        "setAttr": 3,
        "delAttr": 4,
        "setText": 5,
        "appendChild": 6,
        "removeChild": 7,
        "replaceChild": 8,
    })
};

function render(changes = []) {
    changes.forEach(c => {
        switch (c.Action) {
            case goapp.actions.setRoot:
                setRoot(c);
                break;

            case goapp.actions.newNode:
                newNode(c);
                break;

            case goapp.actions.delNode:
                delNode(c);
                break;

            case goapp.actions.setAttr:
                setAttr(c);
                break;

            case goapp.actions.delAttr:
                delAttr(c);
                break;

            case goapp.actions.setText:
                setText(c);
                break;

            case goapp.actions.appendChild:
                appendChild(c);
                break;

            case goapp.actions.removeChild:
                removeChild(c);
                break;

            case goapp.actions.replaceChild:
                replaceChild(c);
                break;

            default:
                console.log(c.Type + ' change is not supported');
        }
    });
}

function setRoot(change = {}) {
    const { NodeID } = change;

    const n = goapp.nodes[NodeID];
    const root = compoRoot(n);

    if (!root) {
        return;
    }

    document.body.replaceChild(root, document.body.firstChild());
}

function newNode(change = {}) {
    const { IsCompo, Type, NodeID, Namespace } = change;

    if (IsCompo) {
        goapp.nodes[NodeID] = {
            Type: Type,
            ID: NodeID,
            IsCompo: true
        };

        return;
    }

    var n = null

    if (Type === 'text') {
        n = document.createTextNode("");
    } else if (change.Namespace) {
        n = document.createElementNS(Namespace, Type);
    } else {
        n = document.createElement(Type);
    }

    n.ID = NodeID;
    goapp.nodes[NodeID] = n;
}

function delNode(change = {}) {
    const { NodeID } = change;
    delete goapp.nodes[NodeID];
}

function setAttr(change = {}) {
    const { NodeID, Key, Value = '' } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    n.setAttribute(Key, Value);
}

function delAttr(change = {}) {
    const { NodeID, Key } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    n.removeAttribute(Key);
}

function setText(change = {}) {
    const { NodeID, Value } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    n.nodeValue = Value;
}

function appendChild(change = {}) {
    const { NodeID, ChildID } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    const c = compoRoot(goapp.nodes[ChildID]);
    if (!c) {
        return;
    }

    if (n.IsCompo) {
        n.RootID = ChildID;
        return;
    }

    n.appendChild(c)
}

function removeChild(change = {}) {
    const { NodeID, ChildID } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    const c = compoRoot(goapp.nodes[ChildID]);
    if (!c) {
        return;
    }

    n.removeChild(c);
}

function replaceChild(change = {}) {
    const { NodeID, ChildID, NewChildID } = change;

    const n = goapp.nodes[NodeID];
    if (!n) {
        return;
    }

    const c = compoRoot(goapp.nodes[ChildID]);
    if (!c) {
        return;
    }

    const nc = compoRoot(goapp.nodes[NewChildID]);
    if (!nc) {
        return;
    }

    if (n.IsCompo) {
        n.RootID = NewChildID;
        return;
    }

    n.replaceChild(nc, c);
}

function compoRoot(node) {
    if (!node || !node.IsCompo) {
        return node;
    }

    const n = goapp.nodes[node.RootID];
    return compoRoot(n);
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

function callCompoHandler(elem, event, fieldOrMethod) {
    switch (event.type) {
        case 'change':
            onchangeToGolang(elem, fieldOrMethod);
            break;

        case 'drag':
        case 'dragstart':
        case 'dragend':
        case 'dragexit':
            onDragStartToGolang(elem, event, fieldOrMethod);
            break;

        case 'dragenter':
        case 'dragleave':
        case 'dragover':
        case 'drop':
            ondropToGolang(elem, event, fieldOrMethod);
            break;

        case 'contextmenu':
            event.preventDefault();

        default:
            eventToGolang(elem, event, fieldOrMethod);
            break;
    }
}

function onchangeToGolang(elem, fieldOrMethod) {
    golangRequest(JSON.stringify({
        'CompoID': elem.compoID,
        'FieldOrMethod': fieldOrMethod,
        'JSONValue': JSON.stringify(elem.value)
    }));
}

function onDragStartToGolang(elem, event, fieldOrMethod) {
    const payload = mapObject(event.dataTransfer);
    payload['Data'] = elem.dataset.drag;
    setPayloadSource(payload, elem);

    event.dataTransfer.setData('text', elem.dataset.drag);

    golangRequest(JSON.stringify({
        'CompoID': elem.compoID,
        'FieldOrMethod': fieldOrMethod,
        'JSONValue': JSON.stringify(payload)
    }));
}

function ondropToGolang(elem, event, fieldOrMethod) {
    event.preventDefault();

    const payload = mapObject(event.dataTransfer);
    payload['Data'] = event.dataTransfer.getData('text');
    payload['FileOverride'] = 'xxx';
    setPayloadSource(payload, elem);

    golangRequest(JSON.stringify({
        'CompoID': elem.compoID,
        'FieldOrMethod': fieldOrMethod,
        'JSONValue': JSON.stringify(payload),
        'Override': 'Files'
    }));
}

function eventToGolang(elem, event, fieldOrMethod) {
    const payload = mapObject(event);
    setPayloadSource(payload, elem);

    if (elem.contentEditable === 'true') {
        payload['InnerText'] = elem.innerText;
    }

    golangRequest(JSON.stringify({
        'CompoID': elem.compoID,
        'FieldOrMethod': fieldOrMethod,
        'JSONValue': JSON.stringify(payload)
    }));
}

function setPayloadSource(payload, elem) {
    payload['Source'] = {
        'GoappID': elem.ID,
        'CompoID': elem.compoID,
        'ID': elem.id,
        'Class': elem.className,
        'Data': elem.dataset,
        'Value': elem.value
    };
}