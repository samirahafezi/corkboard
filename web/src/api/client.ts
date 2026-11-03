import createClient from 'openapi-fetch'

import type { paths } from './schema.gen'

const baseUrl = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'

export const api = createClient<paths>({ baseUrl })
