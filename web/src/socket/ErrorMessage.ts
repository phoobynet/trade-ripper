import { Message } from './Message'

export interface ErrorMessage extends Message {
  type: 'error'
  msg: string
  count: number
}