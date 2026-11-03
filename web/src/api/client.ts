import createClient from 'openapi-fetch'

import type { paths } from './schema.gen'

const baseUrl = import.meta.env.VITE_API_URL ?? '/api'

export const api = createClient<paths>({ baseUrl })
