import create from 'zustand'

export interface AppStore {
  errorsCount: number
  errorMessages: string[]
  infoMessages: string[]
  restartCount: number
  totalTrades: number
  instrumentClass: string
}

export const useAppStore = create<AppStore>(() => ({
  errorsCount: 0,
  errorMessages: [],
  infoMessages: [],
  restartCount: 0,
  totalTrades: 0,
  instrumentClass: '?',
}))
