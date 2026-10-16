const assert = require("node:assert/strict");
const vm = require("node:vm");
const source = process.argv[2];

function deferred() {
  let resolve, reject;
  const promise = new Promise((yes, no) => { resolve = yes; reject = no; });
  return { promise, resolve, reject };
}

function worker(storage = new Map(), hooks = {}) {
  const listeners = {};
  const calls = [];
  const context = {
    console: { log() {}, error() {} },
    self: {
      addEventListener(type, listener) { listeners[type] = listener; },
      async skipWaiting() { calls.push("skipWaiting"); },
      clients: { async claim() { calls.push("claim"); } },
    },
    caches: {
      async open(name) {
        if (!storage.has(name)) storage.set(name, new Map());
        return {
          async addAll(resources) {
            calls.push("addAll");
            assert.ok(resources.includes("/web/app.wasm"));
            if (hooks.populate) await hooks.populate;
            storage.get(name).set("/shared.js", "current version");
            calls.push("populated");
          },
        };
      },
      async keys() { return [...storage.keys()]; },
      async delete(name) {
        if (hooks.cleanup) await hooks.cleanup;
        calls.push("delete");
        return storage.delete(name);
      },
      async match(request, options) {
        const names = options?.cacheName ? [options.cacheName] : [...storage.keys()];
        for (const name of names) {
          const response = storage.get(name)?.get(request);
          if (response !== undefined) return response;
        }
      },
    },
    async fetch(request) { return "network:" + request; },
  };
  vm.runInNewContext(source, context);
  return {
    calls,
    lifecycle(type) {
      const promises = [];
      let dispatching = true;
      listeners[type]({ waitUntil(promise) {
        assert.ok(dispatching, "waitUntil must be called during event dispatch");
        promises.push(promise);
      } });
      dispatching = false;
      assert.equal(promises.length, 1, `${type} must register its work`);
      return Promise.all(promises);
    },
    fetch(request) {
      let response;
      listeners.fetch({ request, respondWith(promise) { response = promise; } });
      return response;
    },
  };
}

async function main() {
  const population = deferred();
  const delayed = worker(new Map(), { populate: population.promise });
  let installed = false;
  const install = delayed.lifecycle("install").then(() => { installed = true; });
  await new Promise(setImmediate);
  assert.equal(installed, false);
  assert.deepEqual(delayed.calls, ["addAll"]);
  population.resolve();
  await install;
  assert.deepEqual(delayed.calls, ["addAll", "populated", "skipWaiting"]);

  const failure = deferred();
  const failed = worker(new Map(), { populate: failure.promise });
  const rejected = assert.rejects(failed.lifecycle("install"), /download failed/);
  failure.reject(new Error("download failed"));
  await rejected;
  assert.deepEqual(failed.calls, ["addAll"]);

  // Old caches coexist with the new cache until activation finishes.
  const storage = new Map([
    ["app-v0", new Map()],
    ["app-v1", new Map([
      ["/shared.js", "old version"],
      ["/only-old", "old version"],
    ])],
  ]);
  const cleanup = deferred();
  const a = worker(storage, { cleanup: cleanup.promise });
  await a.lifecycle("install");
  assert.equal(await a.fetch("/shared.js"), "current version");
  assert.equal(await a.fetch("/only-old"), "network:/only-old");

  let activated = false;
  const activation = a.lifecycle("activate").then(() => { activated = true; });
  await new Promise(setImmediate);
  assert.equal(activated, false);
  assert.ok(!a.calls.includes("claim"));
  cleanup.resolve();
  await activation;
  assert.deepEqual(a.calls.slice(-3), ["delete", "delete", "claim"]);
  assert.deepEqual([...storage.keys()], ["app-v2"]);
  assert.equal(await a.fetch("/shared.js"), "current version");
  assert.equal(await a.fetch("/missing"), "network:/missing");

  // Cleanup is best effort, as before; a failed deletion must not prevent claim.
  const cleanupFailure = deferred();
  const c = worker(new Map([["app-v1", new Map()]]), {
    cleanup: cleanupFailure.promise,
  });
  const bestEffort = c.lifecycle("activate");
  cleanupFailure.reject(new Error("delete failed"));
  await bestEffort;
  assert.deepEqual(c.calls, ["claim"]);
  console.log("Service-worker lifecycle and versioned caches passed");
}

main().catch(error => { console.error(error); process.exitCode = 1; });
