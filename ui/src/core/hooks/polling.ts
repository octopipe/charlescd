import { useCallback, useEffect, useState } from "react"
import axios, { AxiosResponse } from "axios"
import { useAppDispatch } from "./redux"
import { AsyncThunk } from "@reduxjs/toolkit"

interface FetchOptions<T> {
  timer?: any
  thunk: any
}

interface Response<T> {
  loading: boolean
  startPolling: () => Promise<void>
  stopPolling: () => void
}


const usePolling = <T>(options?: FetchOptions<T>): Response<T> => {
  const [error,setError] = useState<Error>()
  const [loading,setLoading] = useState(false)
  const dispatch = useAppDispatch()
  let interval: any
  

  const doRequest = async() => {
    try{
      await dispatch(options?.thunk)
    }catch(err){
      setError(err as Error)
    }
  }
    
  const startPolling = useCallback(async() => {
    setLoading(true)
    await doRequest()
    interval = setInterval(doRequest, options?.timer)
    setLoading(false)
  }, [])

  const stopPolling = useCallback(() => {
    clearInterval(interval)
  }, [])

  return { startPolling, stopPolling, loading }
}

export default usePolling