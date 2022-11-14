import { Gapper } from '../types/Gapper'
import numeral from 'numeral'

type Props = {
  title?: string
  gappers?: Gapper[]
}

export default function Gappers({ gappers, title }: Props) {
  if (!gappers) {
    return null
  }

  return (
    <div className="w-full">
      <header>
        <h3>{title}</h3>
      </header>
      <div className="border border-slate-700 p-2">
        <table className="w-full text-xs">
          <thead className="bg-black text-white">
            <tr>
              <th className="text-left">Ticker</th>
              <th className="text-right">Price</th>
              <th className="text-right">Previous Close</th>
              <th className="text-right">Change</th>
              <th className="text-right">Change %</th>
            </tr>
          </thead>
          <tbody>
            {gappers.map((gapper) => (
              <tr key={gapper.ticker}>
                <td className="font-bold tracking-wider">{gapper.ticker}</td>
                <td className="text-right tabular-nums">
                  {numeral(gapper.p).format('$0,0.00')}
                </td>
                <td className="text-right tabular-nums">
                  {numeral(gapper.pc).format('$0,0.00')}
                </td>
                <td className="text-right tabular-nums">
                  {numeral(gapper.c).format('0,0.00')}
                </td>
                <td className="text-right tabular-nums">
                  {numeral(gapper.cp).format('0,0.00%')}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
