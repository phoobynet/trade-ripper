import { ReactNode, useMemo } from 'react'

type Props = {
  title: string
  comment?: string
  value: string | ReactNode
  type?: 'info' | 'warning' | 'error'
}

export default function Stat({ title, value, type, comment }: Props) {
  const textAccent = useMemo<string>(() => {
    switch (type) {
      case 'info':
        return 'text-info'
      case 'warning':
        return 'text-warning'
      case 'error':
        return 'text-error'
      default:
        return 'text-primary-content'
    }
  }, [type])

  const borderAccent = useMemo<string>(() => {
    switch (type) {
      case 'info':
        return 'border-info'
      case 'warning':
        return 'border-warning'
      case 'error':
        return 'border-error'
      default:
        return 'border-slate-500'
    }
  }, [type])

  return (
    <div
      className={`${borderAccent} border border-slate-500 p-2 rounded-md h-[120px] flex flex-col space-y-2`}
    >
      <header>
        <div
          className={`${textAccent} text-center uppercase tracking-wider font-bold text-xl`}
        >
          {title}
        </div>
      </header>
      <main className={'text-center'}>
        <div className={`text-4xl font-bold tabular-nums text-primary-content`}>
          {value}
        </div>
        <div>{comment ?? <p>{comment}</p>}</div>
      </main>
    </div>
  )
}
