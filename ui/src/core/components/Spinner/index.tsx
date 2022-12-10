import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { useEffect, useState } from 'react'
import { default as SpinnerComponent } from 'react-bootstrap/Spinner';
import './style.scss'


interface Props {

}


const Spinner = (props: Props) => {
  return (
    <SpinnerComponent className='spinner' animation="border" role="status">
      <span className="visually-hidden">Loading...</span>
    </SpinnerComponent>
  )
}

export default Spinner