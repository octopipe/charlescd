import { IconProp } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import React, { useEffect, useState } from 'react'
import './style.scss'


interface Props {
  icon: IconProp
  iconColor: string
  text: string
  onClick: () => void
}


const FloatingButton = ({ icon, iconColor, text, onClick }: Props) => {
  return (
    <div className='floating-button'>
      <div className='floating-button__btn' onClick={onClick}>
        <FontAwesomeIcon icon={icon} color={iconColor} className="me-1" /> {text}
      </div>
    </div>
  )
}

export default FloatingButton