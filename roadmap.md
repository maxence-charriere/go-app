# Code audit roadmap

Audit date: 2026-10-05. This checklist records the original audit findings and tracks implementation below.

## Scope and evidence

Reviewed the library code and tests in `pkg/app`, `pkg/cache`, `pkg/cli`, `pkg/errors`, `pkg/logs`, and `pkg/analytics`, including rendering generators, generated implementation patterns, native/WASM variants, and the library's browser runtime scripts. `docs/` and `pkg/ui/` were completely excluded. Documentation, documentation examples, and website content were excluded everywhere. Recommendations about the browser scripts and static generation concern library implementation, not website content or presentation.

The checklist is ordered by expected impact, considering severity, affected runtime paths, confidence, and implementation scope. All proposed changes preserve exported signatures and types. Correcting defective behavior can change observable results; each relevant item calls out that compatibility consideration. Scheduling and cache policy changes require particular care.

Validation on Go 1.26.2, darwin/arm64:

- The existing native `pkg/app` suite passed, both normally and with the race detector. Test binaries were compiled into and executed from temporary directories because existing tests create relative files.
- `go test -mod=readonly -race ./pkg/cache ./pkg/cli ./pkg/errors ./pkg/logs ./pkg/analytics` passed.
- `go vet -mod=readonly` passed for all six audited packages.
- The `pkg/app` test binary compiled for `GOOS=js GOARCH=wasm`. Browser/WASM runtime tests were not executed; native tests skip some browser behavior.
- Temporary external programs and Go test overlays reproduced the specific failures described in the checklist. HTML was parsed with `golang.org/x/net/html`; a Node probe exercised the unchanged WASM-download wrapper. Temporary probe files and build artifacts were kept outside the repository.

Passing existing tests did not cover the demonstrated edge cases. Performance numbers below describe the implementation at the time of the audit, are indicative local measurements, and are not claims of achieved improvements. The unused `pkg/cache` package was subsequently removed; its audit-only checklist items have been removed as well.

## Prioritized checklist

- [x] **1. Correct HTML attribute escaping and boolean serialization**

  **Implemented:** Attribute values are HTML-escaped after resource URL resolution. Presence-based attributes omit `false` and minimize `true`; data, ARIA, and enumerated values remain literal, including `hidden="until-found"`. Browser mounting handles boolean presence consistently while retaining live property updates. Exported APIs and intentional raw HTML are unchanged.

  **Validated:** Full native `pkg/app` suite with the race detector, full browser/WASM suite in headless Chrome, and `go vet` passed. Added parsed-output regression tests, server/browser state comparisons, and attribute-value fuzzing (20 seconds, over 187,000 executions).

  **Impact:** Critical.

  **Files:** `pkg/app/node.go:739–752`, `pkg/app/attribute.go:9–41,70–110`; `pkg/app/node_test.go:1008–1085` and generated setters in `pkg/app/gen/html.go`.

  **Found:** Attribute values use `strconv.Quote`, which escapes Go strings. A parsed-output probe of a title containing `hello" autofocus onfocus=alert(1) x="` produced additional `autofocus` and `onfocus` attributes. Newlines and ampersands also fail to round-trip. Separately, all values equal to `"true"` lose their value, while presence-based boolean attributes with `"false"` remain present: `Input().Disabled(false)` emits `disabled="false"`, which still disables the input. Data and enumerated attributes lose literal `"true"` values.

  **Why it matters:** Untrusted text passed through ordinary attribute setters can introduce executable markup. Server rendering also disagrees with browser mounting for valid boolean/data attributes. HTML quoting and presence semantics follow the [HTML attribute syntax](https://html.spec.whatwg.org/multipage/syntax.html#attributes-2).

  **Recommended change:** Resolve resource URLs, then escape attribute contents with HTML escaping inside explicit quotes. Explicitly classify presence-based boolean attributes; omit false values and minimize true values only for that class. Preserve string values for data, ARIA, and enumerated attributes. Retain intentional raw-HTML behavior.

  **Expected benefit:** Removes the demonstrated attribute breakout and preserves attribute values and initial control behavior.

  **Risk/difficulty:** Medium compatibility risk for applications compensating for incorrect output; implementation is localized. Preserve generic `Attr` behavior for legitimate names and values.

  **Verification:** Extend encoding tests with parsed attribute comparisons for quotes, apostrophes, ampersands, newlines, Unicode, and resolved URLs; fuzz values and assert no added attributes/nodes. Test true/false `disabled`, `checked`, and `hidden`, plus data, ARIA, `draggable`, `spellcheck`, and `contenteditable`; compare server output with browser DOM behavior.

- [x] **2. Use the standard signal-context lifecycle**

  **Implemented:** `ContextWithSignals` delegates to `signal.NotifyContext` with the existing signature. The returned cancel function unregisters signal delivery; documentation clarifies that it must be called even after signal or parent cancellation and may restore default signal behavior.

  **Validated:** CLI tests passed with the race detector, along with `go vet ./pkg/cli` and a repository-wide build. Regression tests cover explicit and parent cancellation, already-canceled parents, repeated registration/cancellation, and bounded goroutine counts. Signal delivery, repeated signals, zero-signal arguments, and restoration of default interrupt behavior are tested in isolated subprocesses.

  **Impact:** High.

  **Files:** `pkg/cli/cli.go:125–136`; `pkg/cli/cli_test.go:86–92`.

  **Found:** `ContextWithSignals` never unregisters its signal channel, closes it after one signal, and waits only for a signal. Explicit/parent cancellation leaves its goroutine blocked. A subprocess receiving two signals crashed with `send on closed channel`; creating and canceling 100 contexts retained 100 helper goroutines.

  **Why it matters:** Ordinary shutdown can leak resources or crash a process on subsequent signal delivery.

  **Recommended change:** Implement the existing signature using [signal.NotifyContext](https://pkg.go.dev/os/signal#NotifyContext), returning its stop function as `cancel`, preserving zero-signal and parent-cancellation behavior. Calling the returned cancel/stop releases the registration; signal delivery or parent cancellation alone terminates the helper goroutine.

  **Expected benefit:** Removes a process crash and registration/goroutine leaks while simplifying custom synchronization.

  **Risk/difficulty:** Low. Verify deregistration and restoration of signal handling when cancellation releases resources.

  **Verification:** Subprocess tests for delivery and repeated signals; explicit cancellation, already/later-canceled parents, repeated registration/cancellation, and bounded goroutine counts. Keep real signals isolated from the main test process.

- [x] **3. Synchronize memory-storage reads and iteration**

  **Implemented:** `Contains` uses the existing read lock. `ForEach` snapshots keys under the read lock and invokes callbacks after unlocking, allowing callbacks to modify storage safely. The unused expiration cache and remaining `pkg/cache` API were removed instead of retaining the cache initialization work.

  **Validated:** Full native `pkg/app` race tests, full browser/WASM tests in headless Chrome, repository-wide `go vet`, and a repository-wide build passed. Regression tests cover concurrent access, callbacks that read/set/delete/clear storage, snapshot enumeration, and empty storage. The snapshot regression fails against the original implementation. A local darwin/arm64 benchmark measured zero allocations for empty enumeration and one allocation for nonempty snapshots (160 bytes for 10 keys; 16,384 bytes for 1,000 keys).

  **Impact:** High.

  **Files:** `pkg/app/storage.go:62–115`; `pkg/app/storage_test.go`.

  **Found:** Memory storage `Contains` and `ForEach` access a map without its mutex; a public-API concurrency probe reported a race between `Set` and `ForEach`.

  **Why it matters:** Async application work can race or trigger fatal concurrent-map access.

  **Recommended change:** Lock `Contains`; snapshot enumeration keys under a read lock and invoke `ForEach` callbacks after unlocking.

  **Expected benefit:** Safe concurrent map access without callback deadlocks.

  **Risk/difficulty:** Low. Storage enumeration snapshots allocate proportionally to key count; avoid holding locks around mutating callbacks.

  **Verification:** Race tests mixing storage operations and callbacks that set/delete/clear entries. Include empty-storage behavior and enumeration allocation benchmarks.

- [x] **4. Capture a separate context for each async action handler**

  **Implemented:** `Post` copies the context for each handler and assigns its source to that copy. Async closures and UI dispatches use the handler-specific context without changing registration or scheduling APIs.

  **Validated:** Full native `pkg/app` race tests, full browser/WASM tests in headless Chrome, and `go vet ./pkg/app` passed. Regression tests use real goroutines with delayed and concurrent execution, multiple sources, mixed async/UI handlers, repeated posts, and follow-up dispatches. Running the tests against the original implementation reproduced incorrect sources and a data race.

  **Impact:** High.

  **Files:** `pkg/app/action.go:70–92`; `pkg/app/action_test.go`.

  **Found:** `Post` mutates `ctx.sourceElement` for each handler while async closures capture the same variable. A probe observed the wrong source in 88 of 100 posts; the race detector reported the assignment racing with the async closure read. Existing tests use an inline async fake.

  **Why it matters:** `Src()` and later dispatches can refer to an unrelated component, causing incorrect updates or dropped work after unmounting.

  **Recommended change:** Copy the context inside each loop iteration, assign that handler's source to the copy, and capture the copy.

  **Expected benefit:** Correct handler identity and race-free async dispatch with negligible cost.

  **Risk/difficulty:** Low; exported behavior for correctly associated handlers is preserved.

  **Verification:** Delayed real goroutines with mixed global async/component handlers, multiple sources, repeated posts, and source assertions under `-race`. Retain existing synchronous action tests.

- [ ] **5. Move state callbacks and serialization outside manager locks**

  **Impact:** High.

  **Files:** `pkg/app/state.go:107–112,170–259,287–344,371–425`; state/context tests.

  **Found:** Cleanup invokes `Observer.While` while holding the state mutex. A predicate reading another state reproduced a deadlock. JSON marshal/unmarshal and storage calls also occur under that mutex, permitting custom serialization callbacks to reenter it. Notification enqueueing couples the state lock to queue capacity.

  **Why it matters:** Composable predicates or custom state values can freeze frames; slow serialization and storage extend lock contention.

  **Recommended change:** Snapshot necessary state/observer data under lock; run predicates, serialization, storage, and enqueueing outside it. Reacquire only for commits, using an internal identity/generation check so cleanup cannot delete a newer replacement observer.

  **Expected benefit:** Removes reentrant deadlocks and shortens critical sections.

  **Risk/difficulty:** Medium. Preserve replacement, expiration, persistence, and notification ordering semantics.

  **Verification:** Reentrant `While`, marshal/unmarshal, and storage callbacks; replacement during cleanup; concurrent get/set/cleanup under `-race`; large-observer contention benchmarks. Run the existing state suite in full.

- [ ] **6. Remove Host-dependent local proxy fetching and propagate cancellation**

  **Impact:** High.

  **Files:** `pkg/app/http.go:553–609`; `pkg/app/resource.go`; proxy tests in `pkg/app/http_test.go`.

  **Found:** Local proxy misses construct an upstream URL from incoming `r.Host` and call `http.Get`. A transport probe confirmed that a chosen Host selects the fetch target and that an already-canceled request still fetches. Local resources take a network round trip through the serving application; successful results are then cached by resource path alone.

  **Why it matters:** Request-controlled hosts can select unintended upstreams and influence cached content. Self-fetching wastes network resources, depends on deployment routing/TLS, and continues after the client leaves.

  **Recommended change:** Serve local proxy resources through their configured `http.Handler` using a safely copied request and internal response capture. For remote resolvers, construct a context-bound request while preserving compatible client/transport behavior. Cache only successful results and preserve status/content metadata.

  **Expected benefit:** Eliminates the local Host dependency and loopback round trip; canceled remote work can release resources promptly.

  **Risk/difficulty:** Medium. Check custom resolvers/handlers, compression, error status mapping, and cache behavior. Avoid arbitrary new payload limits that break valid resources.

  **Verification:** Host variation must not change local resource selection. Test canceled/slow remote requests, local resources behind a reverse proxy, custom resolver handlers, compressed responses, and cache hits. Benchmark cold local proxy requests and allocations.

- [x] **7. Serialize HTML without constructing a test engine**

  **Implemented:** `HTMLString` and `PrintHTML` share a byte-encoding helper using the existing node manager and an identity URL resolver. Serialization no longer constructs an engine or installs browser callbacks. `PrintHTML` writes the encoded bytes directly. Existing component roots, rendering behavior, escaping, whitespace, and exported signatures are preserved.

  **Validated:** Full native `pkg/app` suite with the race detector, full browser/WASM suite in headless Chrome, and `go vet ./pkg/app` passed. Added exact-output cases for text, raw HTML, attributes, resource URLs, booleans, nested components, and empty rendering, plus component lifecycle and mounted-root checks. Browser tests verify all six global callback identities and navigation remain intact; the callback regression fails against the original implementation. Local darwin/arm64 benchmarks for a div containing one span measured 112 B/op and 2 allocations for `HTMLString`, and 64 B/op and 1 allocation for `PrintHTML`; 1,000-child benchmarks are included.

  **Impact:** High.

  **Files:** `pkg/app/node.go:79–89`, `pkg/app/testing.go:34–46`, `pkg/app/engine.go:58–72,161–165`, `pkg/app/browser.go:18–25`.

  **Found:** Each `HTMLString` constructs a full engine, including two 4,096-slot queues and storage/router/page state. A small-tree probe measured 67,488 B/op and 32 allocations/op, about 5.9 µs/op. On WASM the construction also installs browser handlers, replacing globals such as `onclick`, `onpopstate`, and `goappNav` with callbacks for the temporary engine.

  **Why it matters:** A serialization utility makes large avoidable allocations and can interfere with navigation in a running application.

  **Recommended change:** Use a local node manager and minimal rendering context with identity URL resolution. Share an internal buffer-rendering helper so `PrintHTML` can write the bytes directly.

  **Expected benefit:** Removes engine allocation and browser-handler side effects; avoids the string-to-byte copy in `PrintHTML`.

  **Risk/difficulty:** Low–medium. Preserve component rendering, URL resolution, whitespace, and public signatures.

  **Verification:** Existing output tests; `HTMLString`/`PrintHTML` benchmarks on small and large trees with `-benchmem`. Browser tests must show active callback identities and navigation survive serialization.

- [ ] **8. Fix service-worker lifecycle and cache ownership**

  **Impact:** High.

  **Files:** `pkg/app/gen/app-worker.js:4–57`; generated `pkg/app/scripts.go`; worker tests in `pkg/app/http_test.go`.

  **Found:** Async install/activate listeners never register their promises with `event.waitUntil`, and install errors are swallowed. Cleanup deletes every origin cache except the current version; fetch searches all origin caches. Cache names identify only version, not registration/application.

  **Why it matters:** Installation can advance without completed resource caching, lifecycle work can be interrupted, and upgrades can delete unrelated application caches or retrieve another cache's response. The [service-worker event lifetime rules](https://w3c.github.io/ServiceWorker/#extendableevent-waituntil-method) require lifecycle promises to be registered.

  **Recommended change:** Register lifecycle work synchronously with `waitUntil`; propagate installation failure after logging. Derive a registration/application-scoped namespace, delete only caches with proven ownership, and fetch from the current app cache. Handle legacy names conservatively. Update the generator input and regenerate only during implementation.

  **Expected benefit:** Reliable offline upgrades and isolation between applications sharing an origin.

  **Risk/difficulty:** Medium. Requires an upgrade/migration strategy and browser integration tests; preserve intended skip-waiting/claim ordering and custom template support.

  **Verification:** Delayed/failing population, failed upgrade preserving the old worker, offline reload, activation completion, two registrations plus an unrelated cache, and distinct cached responses for the same URL. Test generated template consistency.

- [ ] **9. Make child filtering linear in input and emitted children**

  **Impact:** Medium.

  **Files:** `pkg/app/node.go:37–74`, `pkg/app/range.go`, generated `Body` methods; `pkg/app/node_test.go:102`.

  **Found:** Each nil removal and selector expansion copies a suffix. Alternating nil/retained inputs measured approximately 0.626 µs at 100 children, 33.8 µs at 1,000, and 6.05 ms at 10,000, including a slice copy per iteration. The existing benchmark has one top-level element.

  **Why it matters:** Condition-heavy or large lists repeatedly incur quadratic copying during construction/rendering.

  **Recommended change:** Use stable compaction/flattening with O(input + emitted children) work and safe expansion handling. Preserve child order, selector behavior, typed-nil filtering, and observable slice-aliasing behavior; clear unused reference slots.

  **Expected benefit:** Less copying and frame-time growth for large lists, with simpler filtering logic.

  **Risk/difficulty:** Medium. Overlapping input/output and selector children that alias input need careful handling.

  **Verification:** Existing filtering tests, empty/nil/typed-nil/nested-selector cases, aliasing/order tests, and fuzz comparisons against current valid outputs. Benchmark flat, nil-heavy, and expanding-selector inputs across increasing sizes with allocations.

- [ ] **10. Preserve status, backpressure, and cancellation in WASM downloads**

  **Impact:** Medium.

  **Files:** `pkg/app/gen/app.js:265–316`; generated `pkg/app/scripts.go`.

  **Found:** Response status options are passed to `ReadableStream` instead of `Response`. A probe converted upstream 404/“Not Found” into 200/empty status text. Its eager `start` loop read all 100 chunks plus EOF before any downstream read, without cancellation forwarding. Unknown lengths produce NaN progress.

  **Why it matters:** Error responses lose status and slow/canceled consumers can leave an entire download buffered.

  **Recommended change:** Pass metadata to the [Response constructor](https://fetch.spec.whatwg.org/#dom-response). Use a pull-based source or compatible transform honoring [stream demand and cancellation](https://streams.spec.whatwg.org/#rs-constructor); forward stream errors and handle unknown lengths deliberately.

  **Expected benefit:** Correct loader failures and reduced queued memory relative to consumer demand.

  **Risk/difficulty:** Medium. Preserve supported browser behavior, streaming compilation, headers, and progress callbacks.

  **Verification:** JS tests for 200/404, invalid/missing lengths, read failure, cancellation, and slow/absent consumers. Browser startup and peak-memory comparisons for large WASM binaries; verify generated script consistency.

- [ ] **11. Encode JavaScript configuration as string literals**

  **Impact:** Medium.

  **Files:** `pkg/app/http.go:322–329,362–365`, `pkg/app/gen/app.js:16–19`, `pkg/app/gen/app-worker.js:4,9`; generated `pkg/app/scripts.go` and HTTP tests.

  **Found:** Configuration text is substituted directly inside quoted JavaScript strings. A public-API probe with a loading label containing quotes and a newline produced an `app.js` file that failed `node --check` with a syntax error. Resolver URLs and other string configuration use the same substitution pattern; worker version text is also inserted into quoted strings.

  **Why it matters:** Valid label text can prevent application initialization, and raw substitution makes output depend on JavaScript syntax characters in configuration.

  **Recommended change:** Generate built-in scripts with JSON-encoded string literals and matching template placeholders, reusing the existing JSON serialization approach. Preserve resource resolution and the public custom-worker-template placeholder contract; distinguish built-in structured generation from caller-provided template behavior.

  **Expected benefit:** Valid configuration strings round-trip into executable JavaScript reliably.

  **Risk/difficulty:** Low–medium. Template quoting and generated/source consistency need coordinated changes; avoid altering custom-template expectations inadvertently.

  **Verification:** Generate and parse scripts using quotes, backslashes, newlines, Unicode, and resolver URLs. Assert decoded configuration values, retain ordinary-output/resource tests, and exercise application/worker startup in browser fixtures.

- [ ] **12. Bind rendering and delayed work to request cancellation**

  **Impact:** Medium.

  **Files:** `pkg/app/http.go:618–639`, `pkg/app/context.go:226–230`, `pkg/app/engine.go:293–318,353–366`.

  **Found:** Page rendering starts from `context.Background`, discarding request cancellation and values. `Context.After` always sleeps for its full duration then dispatches; a canceled-context probe still invoked its callback. Server consumption waits unconditionally for tracked work.

  **Why it matters:** Disconnected clients and canceled work retain renderer/component resources and continue consuming time.

  **Recommended change:** Derive rendering from `r.Context()` with explicit render lifetime cleanup. Make delayed work select between a timer and cancellation; make internal enqueue/consumption cancellation-aware. Ensure async bookkeeping completes on every exit. Arbitrary user functions still need to honor their context; they cannot be forcibly stopped.

  **Expected benefit:** Earlier release of pending timers and request-bound work, and preserved request context values.

  **Risk/difficulty:** Medium. Preserve noncanceled SSR completion and scheduling semantics; coordinate with the queue fix.

  **Verification:** Already/mid-wait canceled timers, normal timers, canceled HTTP render contexts and propagated values, pending-worker completion, and bounded goroutine lifetime. Run engine/context/HTTP tests under `-race`.

- [ ] **13. Resolve route factories under lock and invoke them after unlocking**

  **Impact:** Medium.

  **Files:** `pkg/app/route.go:54–65`; route enumeration in `pkg/app/static.go:34–36`; route/static tests.

  **Found:** Exact and regexp factories execute under the router read lock. A factory registering another route reproduced a deadlock. Static generation directly ranges over `routes.routes`, bypassing router synchronization.

  **Why it matters:** Factories cannot safely compose route registration, and enumeration can race with registration.

  **Recommended change:** Select/copy the factory while locked, then unlock before calling it. Snapshot route paths under the router lock for internal enumeration. Preserve exact-route precedence and regexp insertion order.

  **Expected benefit:** Removes reentrant deadlocks and unsynchronized map enumeration with small internal changes.

  **Risk/difficulty:** Low. Maintain the lookup snapshot's meaning when registration changes concurrently.

  **Verification:** Exact/regexp factories registering routes, concurrent registration/lookup/enumeration under `-race`, and existing routing precedence/static generation tests.

- [ ] **14. Copy request URLs correctly when forwarding normalized paths**

  **Impact:** Medium.

  **Files:** `pkg/app/http.go:493–515`; HTTP/resource tests.

  **Found:** Version-prefixed paths are normalized only in a local variable, then the original request is forwarded to the static handler. A probe requesting an existing `/audit/web/audit.txt` with version `audit` received 404. Legacy WASM forwarding shallow-copies the request but modifies its shared `URL`, changing the caller's request path.

  **Why it matters:** Supported-looking resource paths fail, and request mutation can surprise surrounding handlers or request reuse.

  **Recommended change:** Copy the URL along with the request and forward the correctly normalized handler path. Reconcile `RawPath` when rewriting escaped paths; preserve query strings and resource resolver semantics.

  **Expected benefit:** Correct versioned static-resource forwarding and isolated request rewriting.

  **Risk/difficulty:** Low–medium. Custom resolvers and absolute/relative local directories require regression coverage.

  **Verification:** Versioned/unversioned static paths, both WASM aliases, escaped paths, queries, and custom/local resolvers. Assert the original request and URL are unchanged after handling.

- [ ] **15. Evaluate conditional requests against the selected representation**

  **Impact:** Medium.

  **Files:** `pkg/app/http.go:481–490,498–538,612–639`; HTTP tests.

  **Found:** One version ETag is checked before routing for every response. A nonexistent path with a matching version returned 304 in a probe. The same early check overrides file-serving validators and bypasses rendering even when page content can change within a version.

  **Why it matters:** Missing or changed representations can be treated as unchanged, and static file validation is disconnected from file content.

  **Recommended change:** Select the resource first. Keep compatible version validation for immutable generated assets, let the static handler validate its files, and use a representation-valid strategy for dynamic pages. Handle conditional methods and ETag lists/weak comparisons correctly; preserve the `Version` field's existing role.

  **Expected benefit:** Correct 404/conditional behavior and reduced stale content.

  **Risk/difficulty:** Medium. Validator changes affect observable caching; avoid caching request-dependent pages by URL alone.

  **Verification:** Unknown routes with validators, changed local files, page changes within one version, GET/HEAD and other methods, ETag lists/weak tags, generated assets, and legacy aliases. Measure any extra dynamic-page work.

- [ ] **16. Retain passive options in mounted event metadata**

  **Impact:** Medium.

  **Files:** `pkg/app/node.go:190–214,484–508`, `pkg/app/event.go:92–105`; node/event tests.

  **Found:** Mounted handlers omit the original `passive` field although equality compares it. Unchanged passive renders therefore replace listeners; a passive-to-nonpassive transition can compare equal and retain the old browser option.

  **Why it matters:** Repeated updates create unnecessary JS callbacks/listener churn and can leave `PreventDefault` ineffective after an option change.

  **Recommended change:** Preserve the original handler metadata when attaching JS callback and cleanup fields.

  **Expected benefit:** Stable unchanged listeners and correct option transitions.

  **Risk/difficulty:** Low.

  **Verification:** Both passive transitions and unchanged passive updates; browser listener add/remove counts and cancelable wheel/touch behavior. Compare update allocations and retain existing option/equality tests.

- [ ] **17. Bound file lifetimes and reject failed static-generation responses**

  **Impact:** Medium.

  **Files:** `pkg/app/static.go:55–92,98–141`; `pkg/app/static_test.go`.

  **Found:** Each output file is closed with a defer inside the outer loop, so all descriptors stay open until generation returns. Close errors are ignored. `createStaticPage` accepts error responses as successful output; existing tests check file existence even for pages without registered routes.

  **Why it matters:** Large route sets can exhaust descriptors, and failed renders can produce misleading successful output files.

  **Recommended change:** Scope each file operation so it closes before the next iteration, propagate write/close failures, and validate response status. Prefer obtaining valid content before replacing output, or an internal temporary-file/rename strategy when needed. Make directory creation propagate real filesystem failures.

  **Expected benefit:** Bounded descriptor usage and reliable generator error reporting.

  **Risk/difficulty:** Medium. Returning errors for previously accepted 404/500 output changes defective behavior; keep successful layout/path semantics intact.

  **Verification:** Many-route generation with bounded descriptor usage, injected 404/500/write/close/directory failures, and expected content assertions. Update existing fixtures to register valid pages. Benchmark large route sets and peak memory/descriptors.

- [ ] **18. Validate reflective inputs and propagate CLI decoding failures**

  **Impact:** Medium.

  **Files:** `pkg/cli/option.go:27–38,52–90`, `pkg/cli/cli.go:80–96`, `pkg/app/state.go:450–472`; option/CLI/state tests.

  **Found:** CLI environment decoding ignores `Set` errors; an invalid integer retained its default and reported success. Typed-nil options panic while constructing a reflection validation error. State assignment similarly dereferences typed-nil sources/receivers without checking validity. CLI error output also suppresses most ordinary parse failures.

  **Why it matters:** Invalid configuration silently changes process behavior, and optional/malformed values bypass intended error handling with runtime panics.

  **Recommended change:** Check validity, nilness, and receiver assignability before reflective operations. Propagate option decoding errors through the existing error path with option/environment-key context, excluding sensitive values. Preserve valid flag-over-environment precedence and decide explicitly how an overriding valid flag resolves an invalid environment value. Print ordinary parse failures through the existing reporting mechanism.

  **Expected benefit:** Predictable validation and actionable errors without exported API changes.

  **Risk/difficulty:** Medium. Rejecting malformed inputs changes current behavior; preserve all supported input/copy/precedence semantics.

  **Verification:** Nil interface/typed-nil source/receiver tables, valid pointer/value state copies, mismatches, and invalid scalar/duration/time/collection/nested environment values. Use isolated environment settings and test overrides and error output.

- [ ] **19. Release promise bridge callbacks on every settlement path**

  **Impact:** Medium.

  **Files:** `pkg/app/js_wasm.go:67–82`; WASM JS-bridge tests.

  **Found:** `Value.Then` releases its Go-backed JS callback only after fulfillment calls the user function. Rejection leaves the registered callback retained; a user-callback panic bypasses the subsequent release. This finding is based on code inspection, not a browser heap measurement.

  **Why it matters:** Repeated failed async operations can accumulate retained Go/JS callback resources.

  **Recommended change:** Add internal settlement cleanup for fulfillment and rejection, with deferred cleanup around user callbacks. Preserve fulfillment-only callback semantics and existing rejection reporting; avoid a public promise API redesign.

  **Expected benefit:** Bounded callback lifetime for settled promises and exception paths.

  **Risk/difficulty:** Medium. Cleanup must run once and must not release a function before pending invocation or silently consume rejection semantics.

  **Verification:** Browser/WASM tests for fulfillment, rejection, callback failure, and many rejected promises; instrument callback registrations/releases and compare retained resources after settlement. Native stub tests cannot establish this improvement.

- [ ] **20. Make performance benchmarks exercise valid, bounded workloads**

  **Impact:** Medium.

  **Files:** `pkg/app/http_test.go:590–609`, `pkg/app/node_test.go:102`, `pkg/app/html_test.go`; relevant cache/serialization benchmark coverage.

  **Found:** Handler benchmarks request `/hello`, which is not registered in the global test router, so they measure 404 handling. They reuse one response recorder without resetting its accumulated body; the cold benchmark also makes two requests per iteration. Filtering benchmarks omit the large/nil/selector workloads that expose quadratic cost.

  **Why it matters:** These measurements cannot reliably guide the runtime optimizations above and include workload artifacts.

  **Recommended change:** Use registered routes and validate expected behavior outside timed sections. Reset or replace bounded writers, define one consistent operation per iteration, and separate initialization, rendering, cache-hit, and cold-proxy cases. Add targeted serialization, filter-scaling, replacement, and parallel-read benchmarks where relevant to selected fixes.

  **Expected benefit:** Reproducible performance/allocation evidence for implementation decisions.

  **Risk/difficulty:** Low. Preserve router/global-state isolation so benchmark order cannot affect results.

  **Verification:** Assert route/status/body expectations in setup, report allocations, use representative input sizes, and compare repeated runs on the same toolchain/machine. Ensure memory remains bounded as iteration count grows; use profiles to attribute changes.

- [ ] **21. Traverse modern multi-error trees for metadata lookup**

  **Impact:** Low.

  **Files:** `pkg/errors/errors.go:87–96,101–111,121–131`; error tests.

  **Found:** `HasType`, `Tag`, and `UIError` follow only `Unwrap() error`. A joined plain/enriched error matched through this package's standard-library-backed `Is`, but metadata lookup missed the enriched branch and UI output fell back to combined internal text.

  **Why it matters:** Applications adopting `errors.Join` or multiple `%w` operands lose metadata supported on ordinary wrapping chains.

  **Recommended change:** Add an internal ordered traversal supporting both unwrap forms. Preserve existing single-chain precedence/fallbacks and public `Unwrap` semantics, including tested nil behavior.

  **Expected benefit:** Consistent metadata retrieval with modern standard-library error composition.

  **Risk/difficulty:** Low–medium. Define first-match order across branches without changing existing linear results.

  **Verification:** Nested joins, multiple `%w`, mixed wrappers, competing matches, nil children, and no-match cases; retain the full existing suite.

- [ ] **22. Reset text-node mounted state during dismount**

  **Impact:** Low.

  **Files:** `pkg/app/node.go:284–296`, `pkg/app/text.go:24–29`; node/text lifecycle tests.

  **Found:** The text branch of `Dismount` is empty, leaving its JS value attached and `Mounted()` true, unlike HTML/raw nodes. Subsequent mounting rejects an already-mounted text node. This is a direct code-path finding; no browser retention measurement was performed.

  **Why it matters:** Reused text nodes retain stale mounted state and references after removal.

  **Recommended change:** Clear the text node's mounted JS reference consistently during dismount, without redesigning parent relationships or node ownership.

  **Expected benefit:** Correct lifecycle state and remount behavior; release of the stale JS reference.

  **Risk/difficulty:** Low. Verify replacement/reuse behavior and avoid changes to unrelated lifecycle callbacks.

  **Verification:** Root and nested text mount/dismount/remount tests, text replacement, `Mounted()` assertions, and a browser check of removed-node references.

## Implementation boundaries

Select checklist items before implementation. Generated fluent HTML APIs should retain their public surface; fix their shared runtime or generator inputs when appropriate. Keep current route precedence, scheduling order, and supported resolver behavior unless a separately reviewed bug fix requires a narrowly scoped behavioral correction.

The review of `pkg/logs` and `pkg/analytics` did not establish another material, compatibility-preserving first-pass change. Broad rewrites of exported mutable globals, blanket pooling, replacement of reflection, and mass edits of generated methods would need workload or contract evidence beyond this audit. The unused private URL predicates in `pkg/app/http.go` are cleanup candidates, but their removal alone offers little material benefit and is not prioritized.
