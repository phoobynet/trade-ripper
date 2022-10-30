import { useAppStore } from '../stores/useAppStore'
import { LogMessage } from './Message'
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
      try {
        const logMessage = JSON.parse(ev.data as string) as LogMessage
        useAppStore.setState((state) => ({
          logMessages: takeRight([...state.logMessages, logMessage], 100),
        }))
      } catch (e) {
        console.error('Error parsing message:', ev.data)
      }
    }
  })
}
