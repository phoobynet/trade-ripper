import TopBar from './components/TopBar'
import Dashboard from './routes/Dashboard'
import { Route, Routes } from 'react-router-dom'

export default function App() {
  return (
    <>
      <TopBar />
      <div className={'container mx-auto max-w-4xl pt-4'}>
        <Routes>
          <Route
            index
            element={<Dashboard />}
          ></Route>
        </Routes>
      </div>
    </>
  )
}
