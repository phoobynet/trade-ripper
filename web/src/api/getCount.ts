import { InstrumentClass } from '../types/InstrumentClass'
import { getData } from './http'

export interface ClassResult {
  class: InstrumentClass
}

export const getClass = async (): Promise<InstrumentClass> => {
  return getData<ClassResult>('/class').then((c) => c.class)
}
