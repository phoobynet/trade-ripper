import Stat from '../components/Stat'
import { useAppStore } from '../stores/useAppStore'
import { differenceInSeconds } from 'date-fns'
import { motion } from 'framer-motion'
import numeral from 'numeral'
import { useEffect, useMemo, useRef, useState } from 'react'
import { AiFillAlert, AiOutlineStock } from 'react-icons/ai'
import { BsCurrencyBitcoin } from 'react-icons/bs'

export default function Dashboard() {
  const count = useAppStore((state) => state.count)
  const fetchClass = useAppStore((state) => state.fetchClass)
  const instrumentClass = useAppStore((state) => state.instrumentClass)
  const tradesPerSecond = useAppStore((state) => state.tradesPerSecond)
  const [secondsSinceLastCheckIn, setSecondsSinceLastCheckIn] =
    useState<number>(0)
  const lastMessageCheckInterval = useRef<ReturnType<typeof setInterval>>()

  const displayInstrumentClass = useMemo(() => {
    if (instrumentClass === 'us_equity') {
      return 'US Equity'
    } else if (instrumentClass === 'crypto') {
      return 'Crypto'
    }

    return '?'
  }, [instrumentClass])

  useEffect(() => {
    void (async () => {
      await fetchClass()
    })()

    lastMessageCheckInterval.current = setInterval(() => {
      setSecondsSinceLastCheckIn(
        differenceInSeconds(new Date(), useAppStore.getState().lastMessage),
      )
    }, 1000)

    return () => {
      clearInterval(lastMessageCheckInterval.current)
    }
  }, [])

  return (
    <div>
      <main className={'mx-2 flex flex-col space-y-4 md:mx-0'}>
        {secondsSinceLastCheckIn > 5 ? (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.5 }}
          >
            <div className="alert alert-error shadow-lg">
              <div>
                <AiFillAlert
                  size={48}
                  className={'hidden md:block'}
                />

                <AiFillAlert
                  size={64}
                  className={'display md:hidden'}
                />
                <div className={'text-sm md:text-lg'}>
                  <p>
                    No messages received for {secondsSinceLastCheckIn} seconds,
                    please check the server is still running!{' '}
                  </p>
                  <p>Refresh this page after restarting the server.</p>
                </div>
              </div>
            </div>
          </motion.div>
        ) : (
          <></>
        )}
        <section className={'grid grid-cols-1 gap-2 md:grid-cols-2'}>
          <Stat
            title={'Class'}
            value={
              <div className={'flex items-center justify-center space-x-2'}>
                <div>{displayInstrumentClass} </div>
                {instrumentClass === 'crypto' ? (
                  <BsCurrencyBitcoin />
                ) : (
                  <AiOutlineStock />
                )}
              </div>
            }
            type={'info'}
          ></Stat>
          <Stat
            title={'Trades today'}
            value={numeral(count).format('0.000a')}
            type={'info'}
            comment={`Trades per second: ${numeral(tradesPerSecond).format(
              '0,0',
            )}`}
          ></Stat>
        </section>
      </main>
    </div>
  )
}
