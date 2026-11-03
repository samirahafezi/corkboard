import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { render, screen } from '@testing-library/react'
import { beforeEach, describe, expect, it, vi } from 'vitest'

import { api } from './api/client'
import HealthStatus from './HealthStatus'

vi.mock('./api/client', () => ({
  api: { GET: vi.fn() },
}))

const getHealth = vi.mocked(api.GET)

function renderHealthStatus() {
  const queryClient = new QueryClient()

  render(
    <QueryClientProvider client={queryClient}>
      <HealthStatus />
    </QueryClientProvider>,
  )
}

describe('HealthStatus', () => {
  beforeEach(() => {
    getHealth.mockReset()
  })

  it('shows a loading message while the request is pending', () => {
    getHealth.mockReturnValue(new Promise(() => {}))

    renderHealthStatus()

    expect(screen.getByRole('status')).toHaveTextContent('Checking the API…')
  })

  it('shows the status when the API responds', async () => {
    getHealth.mockResolvedValue({
      data: { status: 'ok' },
      response: new Response(),
    })

    renderHealthStatus()

    expect(await screen.findByText('API status: ok')).toBeInTheDocument()
  })

  it('shows unavailable when the API returns no data', async () => {
    getHealth.mockResolvedValue({
      error: undefined,
      response: new Response(null, { status: 502 }),
    } as never)

    renderHealthStatus()

    expect(await screen.findByText('API status: unavailable')).toBeInTheDocument()
  })

  it('shows unavailable when the request fails', async () => {
    getHealth.mockRejectedValue(new TypeError('Failed to fetch'))

    renderHealthStatus()

    expect(await screen.findByText('API status: unavailable')).toBeInTheDocument()
  })
})
