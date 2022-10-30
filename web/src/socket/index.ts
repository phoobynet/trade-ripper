import { AppStore, useAppStore } from '../stores/useAppStore'
import { ErrorMessage } from './ErrorMessage'
import { InfoMessage } from './InfoMessage'
import { RestartMessage } from './RestartMessage'
import { TradeCountMessage } from './TradeCountMessage'
import { takeRight } from 'lodash'

let socket: WebSocket

export const startSocket = () => {
  return new Promise((resolve, reject) => {
    socket = new WebSocket('ws://localhost:3000/ws')
    socket.onopen = () => {
      console.log('Connected to server')
      useAppStore.setState(() => ({
        connectionStatus: 'connected',
        connectionEvent: undefined,
      }))
      resolve(socket)
    }

    socket.onerror = (ev) => {
      console.error('Error connecting to server', ev)
      useAppStore.setState(() => ({
        connectionStatus: 'error',
        connectionEvent: ev,
      }))
      reject(ev)
    }

    socket.onclose = (ev) => {
      console.warn('Connection closed', ev)
      useAppStore.setState(() => ({
        connectionStatus: 'disconnected',
        connectionEvent: ev,
      }))
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
            connectionStatus: state.connectionStatus,
            connectionEvent: undefined,
          }))
        } else if (messageType === 'info') {
          const infoMessage = j as InfoMessage
          useAppStore.setState((state: AppStore) => ({
            infoMessages: takeRight(
              [...state.infoMessages, infoMessage.msg],
              100,
            ),
            connectionStatus: state.connectionStatus,
            connectionEvent: undefined,
          }))
        } else if (messageType === 'restart') {
          const restartMessage = j as RestartMessage
          useAppStore.setState((state) => ({
            restartCount: restartMessage.count,
            connectionStatus: state.connectionStatus,
            connectionEvent: undefined,
          }))
        } else if (messageType === 'tradeCount') {
          const tradeCount = j as TradeCountMessage
          useAppStore.setState((state) => ({
            totalTrades: tradeCount.count,
            connectionStatus: state.connectionStatus,
            connectionEvent: undefined,
          }))
        }
      }
    }
  })
}
