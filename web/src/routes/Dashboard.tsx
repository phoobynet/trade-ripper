import Stat from '../components/Stat'
import { useAppStore } from '../stores/useAppStore'
import numeral from 'numeral'
import { useEffect } from 'react'
import { sentenceCase } from 'sentence-case'

export default function Dashboard() {
  const totalTrades = useAppStore((state) => state.totalTrades)
  const fetchClass = useAppStore((state) => state.fetchClass)
  const instrumentClass = useAppStore((state) => state.instrumentClass)

  useEffect(() => {
    void (async () => {
      await fetchClass()
    })()
  }, [])

  return (
    <div>
      <main className={'flex flex-col space-y-4'}>
        <section className={'grid grid-cols-3 gap-1'}>
          <Stat
            title={'Class'}
            value={sentenceCase(instrumentClass ?? 'loading')}
            type={'info'}
          ></Stat>
          <Stat
            title={'Total Trades'}
            value={numeral(totalTrades).format('0.00a')}
            type={'info'}
          ></Stat>
        </section>
      </main>
    </div>
  )
}
