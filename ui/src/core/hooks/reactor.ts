import { useEffect, useState } from "react"

export interface reactorConfig {
  cb: () => Promise<any>,
  duration: number
}

const useReactor = async (duration: number) => {
  const [res, setRes] = useState<any>()

  const reactor = async (cb: () => Promise<any>) => {
    const r = await cb()
    setRes(r)
  }

  return reactor
}