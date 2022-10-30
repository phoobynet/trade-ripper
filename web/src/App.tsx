import Dashboard from './routes/Dashboard'
import { startSocket } from './socket'
import { useAppStore } from './stores/useAppStore'
import { useEffect, useState } from 'react'
import { Route, Routes } from 'react-router-dom'

export default function App() {
  const [showErrorOverlay, setShowErrorOverlay] = useState(false)
  const connectionStatus = useAppStore((state) => state.connectionStatus)

  useEffect(() => {
    void (async () => {
      try {
        await startSocket()
      } catch (e) {
        setShowErrorOverlay(true)
        console.error(e)
      }
    })()
  }, [])

  useEffect(() => {
    if (connectionStatus === 'disconnected' || connectionStatus === 'error') {
      setShowErrorOverlay(true)
    } else if (connectionStatus === 'connected') {
      setShowErrorOverlay(false)
    }
  }, [connectionStatus])

  return (
    <div className={'container mx-auto max-w-4xl pt-4'}>
      {showErrorOverlay ? (
        <div
          className={
            'absolute top-0 left-0 w-full h-full bg-red-500 text-white'
          }
        >
          <div className={'flex flex-col space-y-10 items-center'}>
            <div className={'text-4xl text-center mt-10'}>
              Houston, we have a problem.
            </div>
            <p>Check that the trade-ripper process is running.</p>
            <div>
              <p className={'text-center'}>
                You can run trade-ripper with the following command
              </p>
              <pre className={'mt-4'}>
                trade-ripper -c <i>[crypto|us_equity]</i> -q{' '}
                <i>my.questdb.host</i>:9009
              </pre>
            </div>
          </div>
        </div>
      ) : (
        <Routes>
          <Route
            index
            element={<Dashboard />}
          ></Route>
        </Routes>
      )}
    </div>
  )
}
