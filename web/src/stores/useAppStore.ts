import { getClass } from '../api/getCount'
import { InstrumentClass } from '../types/InstrumentClass'
import { take } from 'lodash'
import create from 'zustand'

const sourceUrl = import.meta.env.DEV
  ? `http://${location.hostname}:3000`
  : `http://${location.hostname}:${location.port}`

export interface Message {
  type: string
  data: unknown
  time: string
  message: string
}

export interface AppStore {
  lastMessage: Date
  messages: Message[]
  count: number
  rateBuffer: number[]
  tradesPerSecond: number
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
  fetchClass: async () => {
    getClass().then((c) => set({ instrumentClass: c }))
  },
}))

const eventSource = new EventSource(`${sourceUrl}/events?stream=events`)

interface CountMessageData {
  n: number
}

eventSource.onmessage = (event) => {
  const message = JSON.parse(event.data) as Message

  if (message.message === 'count') {
    useAppStore.setState((state) => {
      let count = state.count

      if (message.message === 'count') {
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
  }

  useAppStore.setState((state) => ({
    lastMessage: new Date(),
  }))
}
