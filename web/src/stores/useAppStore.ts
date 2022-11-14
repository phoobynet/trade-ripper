import { getClass } from '../api/getCount'
import { InstrumentClass } from '../types/InstrumentClass'
import { take } from 'lodash'
import create from 'zustand'
import { Gapper } from '../types/Gapper'

const sourceUrl = import.meta.env.DEV
  ? `http://${location.hostname}:3000`
  : `http://${location.hostname}:${location.port}`

export interface Message {
  type: string
  data: unknown
  time?: string
  message: string
}

export interface Calendar {
  date: string
  sessionOpen: string
  sessionClose: string
  open: string
  close: string
}

export interface MarketStatus {
  status: string
  localTime: string
  marketTime: string
  current?: Calendar
  next: Calendar
  previous: Calendar
}

export interface AppStore {
  lastMessage: Date
  messages: Message[]
  count: number
  rateBuffer: number[]
  tradesPerSecond: number
  marketStatus?: MarketStatus
  gappers?: Gapper[]
  instrumentClass?: InstrumentClass
  fetchClass: () => Promise<void>
}

export const useAppStore = create<AppStore>((set) => ({
  lastMessage: new Date(),
  messages: [],
  count: 0,
  rateBuffer: [],
  tradesPerSecond: 0,
  instrumentClass: undefined,
  marketStatus: undefined,
  gappers: undefined,
  fetchClass: async () => {
    getClass().then((c) => set({ instrumentClass: c }))
  },
}))

const eventSource = new EventSource(`${sourceUrl}/api/events?stream=events`)

interface CountMessageData {
  n: number
}

eventSource.onmessage = (event) => {
  const message = JSON.parse(event.data) as Message

  if (message.type === 'trade_count') {
    useAppStore.setState((state) => {
      let count = state.count

      if (message.type === 'trade_count') {
        const countMessageData = message.data as CountMessageData
        count = countMessageData.n
      }

      const rateBuffer = take([count - state.count, ...state.rateBuffer], 60)

      const totalTradesInLastMinute = rateBuffer.reduce((a, b) => a + b, 0)

      return {
        messages: take([message, ...state.messages], 100),
        count,
        rateBuffer,
        tradesPerSecond:
          totalTradesInLastMinute > 0
            ? totalTradesInLastMinute / rateBuffer.length
            : 0,
      }
    })
  } else if (message.type === 'market_status') {
    useAppStore.setState(() => ({
      marketStatus: message.data as MarketStatus,
    }))
  } else if (message.type === 'gappers') {
    useAppStore.setState(() => ({
      gappers: message.data as Gapper[],
    }))
  }

  useAppStore.setState((state) => ({
    lastMessage: new Date(),
  }))
}
