import { useCallback, useEffect, useState } from "react"
import axios from "axios"

const BASE_URL = 'http://localhost:8080'

interface FetchOptions {
  method?: string
  data?: any
}

interface Response<T> {
  data?: T
  error?: Error
  loading: boolean
  fetch: (url: string, options?: FetchOptions) => Promise<any>
}

const useFetch = <T>(): Response<T> => {
  const [data,setData] = useState<T>()
  const [error,setError] = useState<Error>()
  const [loading,setLoading] = useState(false)

  const fetch = useCallback(async(url: string, options?: FetchOptions) => {
    const method = options?.method || 'GET'

    try{
      setLoading(true)
      const response = await axios.request({ url: `${BASE_URL}${url}`, method, data: options?.data})
      setData(response.data)
      return response.data
    }catch(err){
      setError(err as Error)
    }finally{
      setLoading(false)
    }
  }, [])

  return { fetch, error, loading, }
}

export default useFetch