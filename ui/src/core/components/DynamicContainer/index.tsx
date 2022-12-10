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


const DynamicContainer = ({ className, loading = true, children }: Props) => (
  <div className={`dynamic-container ${className}`}>
    {loading ? <Spinner /> : children}
  </div>
)

export default DynamicContainer