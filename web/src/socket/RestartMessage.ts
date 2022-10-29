import { Message } from './Message'

export interface RestartMessage extends Message {
  type: 'restart'
  msg: string
  count: number
}