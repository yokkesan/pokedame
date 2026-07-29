import {
  useEffect,
  useState,
} from 'react'

import './App.css'

type ApiStatus = 'loading' | 'connected' | 'failed'

type HealthResponse = {
  status: string
}

function App() {
  const [apiStatus, setApiStatus] = useState<ApiStatus>('loading')

  useEffect(() => {
    const controller = new AbortController()

    const checkApiConnection = async () => {
      try {
        const response = await fetch('/api/health', {
          method: 'GET',
          headers: {
            Accept: 'application/json',
          },
          signal: controller.signal,
        })

        if (!response.ok) {
          throw new Error(`API responded with ${response.status}`)
        }

        const data = (await response.json()) as HealthResponse

        if (data.status !== 'ok') {
          throw new Error('Unexpected health response')
        }

        setApiStatus('connected')
      } catch (error) {
        if (error instanceof DOMException && error.name === 'AbortError') {
          return
        }

        console.error('API connection failed:', error)
        setApiStatus('failed')
      }
    }

    void checkApiConnection()

    return () => {
      controller.abort()
    }
  }, [])

  return (
    <main className="app">
      <section className="connection-panel">
        <h1>pokedame</h1>

        <p className={`connection-status connection-status--${apiStatus}`}>
          {apiStatus === 'loading' && 'API接続を確認しています'}
          {apiStatus === 'connected' && 'API接続：正常'}
          {apiStatus === 'failed' && 'API接続：失敗'}
        </p>
      </section>
    </main>
  )
}

export default App