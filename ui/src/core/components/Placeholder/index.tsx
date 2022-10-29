import React from "react";
import { ReactSVG } from 'react-svg';
import './style.scss'

interface Props {
  children: JSX.Element
  text: string
}

const Placeholder = React.forwardRef(
  (
    { children, text }: Props,
    ref: React.Ref<HTMLDivElement>
  ) => {
    return (
      <div className="place">
        {children}
        <div className="place__text">{text}</div>
      </div>
    )
  }
)

export default Placeholder