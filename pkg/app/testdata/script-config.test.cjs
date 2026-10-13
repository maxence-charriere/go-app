const assert = require("node:assert/strict");
const fs = require("node:fs");
const vm = require("node:vm");
const test = JSON.parse(fs.readFileSync(0, "utf8"));

async function checkApp() {
  const labels = [];
  const errors = [];
  let wasmURL, workerURL, removed = false, finish;
  const started = new Promise(resolve => { finish = resolve; });
  const context = vm.createContext({
    console: { error(...args) { errors.push(args); } },
    URLSearchParams, ReadableStream, Response,
    window: { location: { search: "" }, addEventListener() {} },
    navigator: {
      userAgent: "script configuration test",
      serviceWorker: {
        async register(url) {
          workerURL = url;
          return { addEventListener() {} };
        },
        addEventListener() {},
      },
    },
    document: {
      getElementById(id) {
        if (id === "app-wasm-loader") return { remove() { removed = true; } };
        if (id === "app-wasm-loader-label") return { set innerText(value) { labels.push(value); } };
        return {};
      },
    },
    Go: class {
      importObject = {};
      run(instance) { assert.equal(instance, "wasm instance"); finish(); }
    },
    WebAssembly: {
      async instantiateStreaming(response) {
        await (await response).arrayBuffer();
        return { instance: "wasm instance" };
      },
    },
    async fetch(url) {
      wasmURL = url;
      return new Response(new Uint8Array([0, 1, 2]), {
        headers: { "Content-Length": "3" },
      });
    },
  });

  // Parsing and executing the complete script also exercises startup calls.
  vm.runInContext(test.app, context);
  await started;
  await new Promise(setImmediate);
  assert.deepEqual(errors, []);
  assert.equal(context.injected, undefined);
  assert.equal(removed, true);
  assert.equal(workerURL, test.workerURL);
  assert.equal(wasmURL, test.wasmURL);
  assert.equal(labels[0], test.value.replace("{progress}", "0"));
  for (const name of ["goappLoadingLabel", "goappWasmContentLength", "goappWasmContentLengthHeader"]) {
    assert.equal(vm.runInContext(name, context), test.value);
  }
  assert.deepEqual(JSON.parse(vm.runInContext("JSON.stringify(goappEnv)", context)), test.env);
}

async function checkWorker() {
  const listeners = {};
  const logs = [];
  const calls = [];
  const context = vm.createContext({
    console: { log(message) { logs.push(message); }, error(error) { throw error; } },
    self: {
      addEventListener(type, listener) { listeners[type] = listener; },
      async skipWaiting() { calls.push("skipWaiting"); },
      clients: { async claim() { calls.push("claim"); } },
    },
    caches: {
      async open(name) {
        assert.equal(name, "app-" + test.value);
        calls.push("open");
        return { async addAll(resources) {
          assert.deepEqual(Array.from(resources).sort(), [...test.resources].sort());
          calls.push("addAll");
        } };
      },
      async keys() { return []; },
      async match(request, options) {
        assert.equal(options.cacheName, "app-" + test.value);
        return "cached response";
      },
    },
  });
  vm.runInContext(test.worker, context);
  assert.equal(context.injected, undefined);
  for (const type of ["install", "activate"]) {
    let work;
    listeners[type]({ waitUntil(promise) { work = promise; } });
    assert.ok(work);
    await work;
  }
  assert.deepEqual(calls, ["open", "addAll", "skipWaiting", "claim"]);
  assert.deepEqual(logs, [
    "installing app worker " + test.value,
    "app worker " + test.value + " is activated",
  ]);
  let response;
  listeners.fetch({ request: "/test", respondWith(promise) { response = promise; } });
  assert.equal(await response, "cached response");
}

(async () => {
  await checkApp();
  await checkWorker();
  console.log("JavaScript configuration checks passed");
})().catch(error => { console.error(error); process.exitCode = 1; });
