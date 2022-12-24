import { useCallback, useEffect, useState } from "react"
import axios, { AxiosResponse } from "axios"

interface FetchOptions<T> {
  timer?: any
  request: () => Promise<AxiosResponse<T>> 
}

interface Response<T> {
  loading: boolean
  data?: T
  error?: Error
  startPolling: () => Promise<any>
  stopPolling: () => void
}


const usePolling = <T>(options?: FetchOptions<T>): Response<T> => {
  const [data, setData] = useState<T>()
  const [error,setError] = useState<Error>()
  const [loading,setLoading] = useState(false)
  let interval: any
  

  const startPolling = useCallback(async() => {
    setLoading(true)
    interval = setInterval(async() => {
      try{
        const response = await options?.request()
        setData(response?.data)
        return response?.data
      }catch(err){
        setError(err as Error)
      } finally {
        setLoading(false)
      }

    }, options?.timer)
  }, [])

  const stopPolling = useCallback(() => {
    clearInterval(interval)
  }, [])

  return { startPolling, stopPolling, error, data, loading }
}

export default usePolling