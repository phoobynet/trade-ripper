import { Message } from './Message'

export interface InfoMessage extends Message {
  type: 'info'
  msg: string
}