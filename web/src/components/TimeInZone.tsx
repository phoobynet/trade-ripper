type Props = {
  time?: string
  label?: string
}

export default function TimeInZone({ time, label }: Props) {
  return (
    <div
      className={
        'flex items-center overflow-hidden rounded rounded-md border border-orange-500 text-xs'
      }
    >
      <div className={'trade-wider px-2 text-slate-300'}>{label}</div>
      <div className={'bg-orange-500 px-2 tabular-nums text-white'}>{time}</div>
    </div>
  )
}
