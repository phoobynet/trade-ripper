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
  totalTrades: number
  instrumentClass?: InstrumentClass
  fetchClass: () => Promise<void>
}

export const useAppStore = create<AppStore>((set) => ({
  lastMessage: new Date(),
  messages: [],
  totalTrades: 0,
  instrumentClass: undefined,
  fetchClass: async () => {
    getClass().then((c) => set({ instrumentClass: c }))
  },
}))

const eventSource = new EventSource(`${sourceUrl}/events?stream=messages`)

eventSource.onmessage = (event) => {
  useAppStore.setState((state) => {
    const message = JSON.parse(event.data) as Message

    let totalTrades = state.totalTrades

    if (message.message === 'tradeCount') {
      totalTrades = (
        message.data as {
          totalTrades: number
        }
      ).totalTrades
    }

    return {
      lastMessage: new Date(),
      messages: take([message, ...state.messages], 100),
      totalTrades,
    }
  })
}
