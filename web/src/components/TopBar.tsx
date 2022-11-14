import { useAppStore } from '../stores/useAppStore'
import { format, parseISO } from 'date-fns'
import { formatInTimeZone } from 'date-fns-tz'
import { get } from 'lodash'
import { useEffect, useState } from 'react'
import TimeInZone from './TimeInZone'
import MarketStatus from './MarketStatus'

const TimeFormat = 'h:mm:ss a'

export default function TopBar() {
  const marketStatus = useAppStore((state) => state.marketStatus)
  const [localTimeFormatted, setLocalTimeFormatted] = useState<string>('')
  const [marketTimeFormatted, setMarketTimeFormatted] = useState<string>('')

  useEffect(() => {
    const lt = get(marketStatus, 'localTime')
    if (lt) {
      setLocalTimeFormatted(format(parseISO(lt), TimeFormat))
    }

    const mt = get(marketStatus, 'marketTime')

    if (mt) {
      setMarketTimeFormatted(
        formatInTimeZone(parseISO(mt), 'America/New_York', TimeFormat),
      )
    }
  }, [marketStatus])
  return (
    <div className={'flex h-10 items-center justify-between px-2'}>
      <div
        className={
          'font-bold uppercase tracking-widest text-slate-200 md:hidden'
        }
      >
        TR
      </div>

      <div
        className={
          'hidden font-bold uppercase tracking-widest text-slate-200 md:block'
        }
      >
        trade ripper
      </div>
      <div className={'flex space-x-2'}>
        <MarketStatus marketStatus={marketStatus?.status} />
        <TimeInZone
          time={localTimeFormatted}
          label={'Local Time'}
        />
        <TimeInZone
          time={marketTimeFormatted}
          label={'New York'}
        />
      </div>
    </div>
  )
}
