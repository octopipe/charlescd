import { useCallback, useEffect, useState } from "react"
import axios from "axios"

const BASE_URL = 'http://localhost:8080'

interface Response<T> {
  data?: T
  error?: Error
  loading: boolean
  fetch: (url: string, method: string, data?: any) => Promise<any>
}

const useFetch = <T>(): Response<T> => {
  const [data,setData] = useState<T>()
  const [error,setError] = useState<Error>()
  const [loading,setLoading] = useState(false)

  const fetch = useCallback(async(url: string, method: string, data?: any) => {
    try{
      setLoading(true)
      const response = await axios.request({ url: `${BASE_URL}${url}`, method, data})
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