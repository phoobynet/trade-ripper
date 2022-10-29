import { Message } from './Message'

export interface TradeCountMessage extends Message {
  type: 'tradeCount'
  count: number
}