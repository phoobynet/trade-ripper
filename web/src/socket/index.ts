import { AppStore, useAppStore } from '../stores/useAppStore'
import { ErrorMessage } from './ErrorMessage'
import { InfoMessage } from './InfoMessage'
import { RestartMessage } from './RestartMessage'
import { TradeCountMessage } from './TradeCountMessage'
import { takeRight } from 'lodash'

const socket = new WebSocket('ws://localhost:3000')

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
    }
  }
}
