import Stat from '../components/Stat'
import { useAppStore } from '../stores/useAppStore'
import numeral from 'numeral'
import { useEffect, useState } from 'react'
import { sentenceCase } from 'sentence-case'

export default function Dashboard() {
  const errorsCount = useAppStore((state) => state.errorsCount)
  const totalTrades = useAppStore((state) => state.totalTrades)
  const infoMessages = useAppStore((state) => state.infoMessages)
  const errorMessages = useAppStore((state) => state.errorMessages)
  const [instrumentClass, setInstrumentClass] = useState<string>('')

  useEffect(() => {
    void (async () => {
      const response = await fetch('http://localhost:3000/api/class')
      setInstrumentClass(
        await response.json().then((j) => (j as { class: string }).class),
      )
    })()
  }, [])

  return (
    <div>
      <main className={'flex flex-col space-y-4'}>
        <section className={'grid grid-cols-3 gap-1'}>
          <Stat
            title={'Class'}
            value={sentenceCase(instrumentClass)}
            type={'info'}
          ></Stat>
          <Stat
            title={'Total Trades'}
            value={numeral(totalTrades).format('0.00a')}
            type={'info'}
          ></Stat>
          <Stat
            title={'Errors'}
            value={numeral(errorsCount).format('0,0')}
            type={'error'}
          ></Stat>
        </section>
        <section>
          <pre>{infoMessages}</pre>
        </section>
        <section>
          <pre>{errorMessages}</pre>
        </section>
      </main>
    </div>
  )
}
