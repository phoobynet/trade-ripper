import Stat from '../components/Stat'
import { useAppStore } from '../stores/useAppStore'
import numeral from 'numeral'
import { sentenceCase } from 'sentence-case'

export default function Dashboard() {
  const errorsCount = useAppStore((state) => state.errorsCount)
  const totalTrades = useAppStore((state) => state.totalTrades)
  const instrumentClass = useAppStore((state) => state.instrumentClass)

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
      </main>
    </div>
  )
}
