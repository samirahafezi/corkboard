// Forwards /api/* requests to the Go API, so the browser only talks to one origin.

interface Env {
  API_ORIGIN?: string
}

interface Context {
  request: Request
  env: Env
}

// apiTarget maps a request URL onto the API origin. Only the path and query
// string are copied, so a request can never change which host is called.
export function apiTarget(requestUrl: string, apiOrigin: string): URL {
  const incoming = new URL(requestUrl)
  const target = new URL(apiOrigin)
  target.pathname = incoming.pathname.replace(/^\/api(?=\/|$)/, '') || '/'
  target.search = incoming.search
  return target
}

export async function onRequest({ request, env }: Context): Promise<Response> {
  if (!env.API_ORIGIN) {
    return new Response('The API proxy is not configured: set API_ORIGIN.', { status: 500 })
  }
  return fetch(new Request(apiTarget(request.url, env.API_ORIGIN), request))
}
