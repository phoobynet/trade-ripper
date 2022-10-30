import { getClass } from '../api/getCount'
import { InstrumentClass } from '../types/InstrumentClass'
import create from 'zustand'

export interface AppStore {
  errorsCount: number
  errorMessages: string[]
  infoMessages: string[]
  restartCount: number
  totalTrades: number
  instrumentClass?: InstrumentClass
  connectionStatus: 'connected' | 'disconnected' | 'error'
  connectionEvent?: Event | CloseEvent
  fetchClass: () => Promise<void>
}

export const useAppStore = create<AppStore>((set) => ({
  errorsCount: 0,
  errorMessages: [],
  infoMessages: [],
  restartCount: 0,
  totalTrades: 0,
  instrumentClass: undefined,
  connectionStatus: 'disconnected',
  connectionEvent: undefined,
  fetchClass: async () => {
    getClass().then((c) => set({ instrumentClass: c }))
  },
}))
