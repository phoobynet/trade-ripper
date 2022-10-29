import { takeRight } from 'lodash'
import create from 'zustand'

// TODO: Make this configuration
const socket = new WebSocket('ws://localhost:3000')

export interface Message {
  type: 'error' | 'info' | 'restart' | 'tradeCount' | 'instrumentClass'
}

export interface InstrumentClassMessage extends Message {
  type: 'instrumentClass'
  instrumentClass: 'us_equity' | 'crypto'
}

export interface ErrorMessage extends Message {
  type: 'error'
  msg: string
  count: number
}

export interface InfoMessage extends Message {
  type: 'info'
  msg: string
}

export interface RestartMessage extends Message {
  type: 'restart'
  msg: string
  count: number
}

export interface TradeCountMessage extends Message {
  type: 'tradeCount'
  count: number
}

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

socket.onopen = () => {
  console.log('Connected to server')
}

socket.onerror = (ev) => {
  console.error('Error connecting to server', ev)
}

socket.onmessage = (ev: MessageEvent) => {
  const j = JSON.parse(ev.data as string)

  if ('type' in j) {
    const messageType = j.type as string
    if (messageType === 'error') {
      const errorMessage = j as ErrorMessage
      useAppStore.setState((state: AppStore) => ({
        errorsCount: errorMessage.count,
        errors: takeRight([...state.errorMessages, errorMessage.msg], 100),
      }))
    } else if (messageType === 'info') {
      const infoMessage = j as InfoMessage
      useAppStore.setState((state: AppStore) => ({
        infoMessages: takeRight([...state.infoMessages, infoMessage.msg], 100),
      }))
    } else if (messageType === 'restart') {
      const restartMessage = j as RestartMessage
      useAppStore.setState(() => ({
        restartCount: restartMessage.count,
      }))
    } else if (messageType === 'tradeCount') {
      const tradeCount = j as TradeCountMessage
      useAppStore.setState(() => ({
        totalTrades: tradeCount.count,
      }))
    } else if (messageType === 'instrumentClass') {
      const instrumentClassMessage = j as InstrumentClassMessage
      useAppStore.setState(() => ({
        instrumentClass: instrumentClassMessage.instrumentClass,
      }))
    }
  }
}
