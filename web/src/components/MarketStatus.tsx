import { useEffect, useState } from 'react'
import { sentenceCase } from 'sentence-case'

type Props = {
  marketStatus?: string
}

export default function MarketStatus({ marketStatus }: Props) {
  const [formattedStatus, setFormattedStatus] = useState<string>('')

  useEffect(() => {
    if (marketStatus) {
      setFormattedStatus(sentenceCase(marketStatus))
    }
  }, [marketStatus])

  return (
    <>
      <div
        className={
          'flex items-center overflow-hidden rounded rounded-md border border-orange-500 text-xs font-bold'
        }
      >
        <div className={'trade-wider px-2 text-slate-300'}>Market Status</div>
        <div className={'bg-orange-500 px-2 tabular-nums text-white'}>
          {formattedStatus}
        </div>
      </div>
    </>
  )
}
