import { useQuery } from '@tanstack/react-query'

import { api } from './api/client'

async function fetchHealth() {
  const { data } = await api.GET('/health')
  if (!data) {
    throw new Error('The API did not return a health status.')
  }
  return data
}

function HealthStatus() {
  const { data, isPending, isError } = useQuery({
    queryKey: ['health'],
    queryFn: fetchHealth,
    retry: false,
  })

  let message = 'Checking the API…'
  if (isError) {
    message = 'API status: broken'
  } else if (!isPending) {
    message = `API status: ${data?.status}`
  }

  return <p role="status">{message}</p>
}

export default HealthStatus
