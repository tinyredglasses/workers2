import "./wasm_exec.js";
import { connect } from 'cloudflare:sockets';

let mod;
let count = 0

globalThis.tryCatch = (fn) => {
  try {
    return {
      result: fn(),
    };
  } catch(e) {
    return {
      error: e,
    };
  }
}

export function init(m) {
  mod = m;
}

async function run(ctx) {
  const go = new Go();

  let ready;
  const readyPromise = new Promise((resolve) => {
    ready = resolve;
  });
  const instance = new WebAssembly.Instance(mod, {
    ...go.importObject,
    workers: {
      ready: () => { ready() }
    },
  });
  go.run(instance, ctx);
  await readyPromise;
}

function createRuntimeContext(env, ctx, binding) {
  return {
    env,
    ctx,
    connect,
    binding,
  };
}

export async function fetch(req, env, ctx) {
  const binding = {};
  await run(createRuntimeContext(env, ctx, binding));
  return binding.handleRequest(req);
}

export async function scheduled(event, env, ctx) {
  const binding = {};
  await run(createRuntimeContext(env, ctx, binding));
  return binding.runScheduler(event);
}

// onRequest handles request to Cloudflare Pages
export async function onRequest(ctx) {
  const binding = {};
  const { request, env } = ctx;
  await run(createRuntimeContext(env, ctx, binding));
  return binding.handleRequest(request);
}

export async function websocketFetch(req, env, ctx) {
  const binding = {};
  await run(createRuntimeContext(env, ctx, binding));


  // const fn = (env, ctx) => binding.handleData(env,ctx)
  try {
    const url = new URL(req.url)
    switch (url.pathname) {
        // case '/':
        // 	return template()
      case '/ws':
        return websocketHandler(req, binding.handleData)
      default:
        return new Response("Not found", { status: 404 })
    }
  } catch (err) {
    return new Response(err.toString())
  }
}

const websocketHandler = async (request, fn) => {
  console.log("websocketHandler")
  const upgradeHeader = request.headers.get("Upgrade")
  if (upgradeHeader !== "websocket") {
    return new Response("Expected websocket", { status: 400 })
  }

  const [client, server] = Object.values(new WebSocketPair())
  await handleWebsocketSession(server, fn)

  return new Response(null, {
    status: 101,
    webSocket: client
  })
}

async function handleWebsocketSession(websocket, fn) {
  console.log("handleSession")
  websocket.accept()
  // await runCode()

  websocket.addEventListener("message", async ({ data }) => {
    // Create instance of WebAssembly Module `mod`, supplying
    // the expected imports in `importObject`. This should be
    // done at the top level of the script to avoid instantiation on every request.

    fn(data)
    if (data === "CLICK") {
      count += 1
      // try {
      // 	const instance = await WebAssembly.instantiate(mod, importObject);
      // 	const retval = instance.exports.Test();
      // 	console.log(retval)
      // } catch (e) {
      // 	console.log(e)
      // }
      websocket.send(JSON.stringify({ count, tz: new Date(), retval: "dfsdf" }))
    } else {
      // An unknown message came into the server. Send back an error message
      websocket.send(JSON.stringify({ error: "Unknown message received", tz: new Date() }))
    }
  })

  websocket.addEventListener("close", async evt => {
    // Handle when a client closes the WebSocket connection
    console.log(evt)
  })
}