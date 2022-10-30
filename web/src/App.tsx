import Dashboard from './routes/Dashboard'
import { startSocket } from './socket'
import { useAppStore } from './stores/useAppStore'
import { useEffect, useRef, useState } from 'react'
import { VscDebugDisconnect } from 'react-icons/vsc'
import { Route, Routes } from 'react-router-dom'

export default function App() {
  const [showErrorOverlay, setShowErrorOverlay] = useState(false)
  const connectionStatus = useAppStore((state) => state.connectionStatus)
  const restartSocketInterval = useRef<
    ReturnType<typeof setInterval> | undefined
  >()

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
      if (!restartSocketInterval.current) {
        restartSocketInterval.current = setInterval(async () => {
          console.info('Attempting to reconnect to socket...')
          await startSocket()
        }, 3_000)
      }
    } else if (connectionStatus === 'connected') {
      clearInterval(restartSocketInterval.current)
      restartSocketInterval.current = undefined
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
          <div className={'flex flex-col space-y-6 items-center'}>
            <div className={'text-4xl text-center mt-10'}>
              Houston, we have a problem.
            </div>
            <div>
              <VscDebugDisconnect size={56}> </VscDebugDisconnect>
            </div>
            <p>Check that the trade-ripper process is running.</p>
            <pre>
              trade-ripper --class <i>[crypto|us_equity]</i> --host{' '}
              <i>my.questdb.host</i> --postgres <i>8812</i> --influx <i>9009</i>
            </pre>
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
