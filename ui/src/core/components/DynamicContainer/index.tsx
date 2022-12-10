import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { Suspense, useEffect, useState } from 'react'
import Spinner from '../Spinner'
import './style.scss'


interface Props {
  className: string
  loading: boolean
  children: React.ReactNode
}


const DynamicContainer = ({ className, loading, children }: Props) => {

  useEffect(() => {
    console.log(loading)

  }, [loading])

  return (
    <div className={`dynamic-container ${className}`}>
      <Suspense fallback={"Loading..."}>
        {children}
      </Suspense>
    </div>
  )
}

export default DynamicContainer