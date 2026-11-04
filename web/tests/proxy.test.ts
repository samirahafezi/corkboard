import { describe, expect, it } from 'vitest'

import { apiTarget } from '../functions/api/[[path]]'

const origin = 'https://corkboard-api.example.com'

describe('apiTarget', () => {
  it.each([
    ['removes the /api prefix', 'https://site.example/api/health', 'https://corkboard-api.example.com/health'],
    ['keeps the query string', 'https://site.example/api/clubs?page=2', 'https://corkboard-api.example.com/clubs?page=2'],
    ['maps /api to the root', 'https://site.example/api', 'https://corkboard-api.example.com/'],
    ['keeps nested paths', 'https://site.example/api/clubs/7/follow', 'https://corkboard-api.example.com/clubs/7/follow'],
    ['only removes a whole /api segment', 'https://site.example/apix/health', 'https://corkboard-api.example.com/apix/health'],
  ])('%s', (_name, requestUrl, want) => {
    expect(apiTarget(requestUrl, origin).href).toBe(want)
  })

  it.each([
    'https://site.example/api//evil.example/steal',
    'https://site.example/api/%2F%2Fevil.example',
    'https://site.example/api/..//evil.example',
    'https://site.example/api/@evil.example',
  ])('never leaves the API host for %s', (requestUrl) => {
    expect(apiTarget(requestUrl, origin).host).toBe('corkboard-api.example.com')
  })
})
