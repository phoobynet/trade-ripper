import Stat from '../components/Stat'
import { useAppStore } from '../stores/useAppStore'
import { differenceInSeconds } from 'date-fns'
import humanizeDuration from 'humanize-duration'
import numeral from 'numeral'
import { useEffect, useRef, useState } from 'react'
import { AiFillAlert } from 'react-icons/ai'
import { sentenceCase } from 'sentence-case'

export default function Dashboard() {
  const totalTrades = useAppStore((state) => state.totalTrades)
  const fetchClass = useAppStore((state) => state.fetchClass)
  const instrumentClass = useAppStore((state) => state.instrumentClass)
  const lastMessage = useAppStore((state) => state.lastMessage)
  const [secondsSinceLastCheckIn, setSecondsSinceLastCheckIn] =
    useState<number>(0)
  const lastMessageCheckInterval = useRef<ReturnType<typeof setInterval>>()

  useEffect(() => {
    void (async () => {
      await fetchClass()
    })()

    lastMessageCheckInterval.current = setInterval(() => {
      setSecondsSinceLastCheckIn(differenceInSeconds(new Date(), lastMessage))
    }, 1000)

    return () => {
      clearInterval(lastMessageCheckInterval.current)
    }
  }, [])

  return (
    <div>
      <main className={'flex flex-col space-y-4 mx-2 md:mx-0'}>
        {secondsSinceLastCheckIn > 5 ? (
          <div className="alert alert-error shadow-lg">
            <div>
              <AiFillAlert
                size={48}
                className={'hidden md:block'}
              />

              <AiFillAlert
                size={64}
                className={'md:hidden display'}
              />
              <span className={'text-sm md:text-lg'}>
                No messages received for{' '}
                {humanizeDuration(secondsSinceLastCheckIn)}, please check the
                server is still running!
              </span>
            </div>
          </div>
        ) : (
          <></>
        )}
        <section className={'grid grid-cols-1 gap-2 md:grid-cols-2'}>
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
