type Props = {
  label: string
  value?: string
}

export default function TopBarLabelValue({ label, value }: Props) {
  return (
    <>
      <div
        className={
          'flex items-center overflow-hidden rounded rounded-md border border-orange-500 text-xs font-bold'
        }
      >
        <div className={'trade-wider px-2 text-slate-300'}>{label}</div>
        <div className={'bg-orange-500 px-2 tabular-nums text-white'}>
          {value}
        </div>
      </div>
    </>
  )
}
