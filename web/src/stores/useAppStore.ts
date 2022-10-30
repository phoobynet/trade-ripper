import { getClass } from '../api/getCount'
import { LogMessage } from '../socket/Message'
import { InstrumentClass } from '../types/InstrumentClass'
import create from 'zustand'

export interface AppStore {
  logMessages: LogMessage[]
  totalTrades: number
  instrumentClass?: InstrumentClass
  connectionStatus: 'connected' | 'disconnected' | 'error'
  connectionEvent?: Event | CloseEvent
  fetchClass: () => Promise<void>
}

export const useAppStore = create<AppStore>((set) => ({
  logMessages: [],
  totalTrades: 0,
  instrumentClass: undefined,
  connectionStatus: 'disconnected',
  connectionEvent: undefined,
  fetchClass: async () => {
    getClass().then((c) => set({ instrumentClass: c }))
  },
}))
