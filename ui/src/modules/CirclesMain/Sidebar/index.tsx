import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useState } from 'react';
import { Link, NavLink, useMatch, useParams, useSearchParams } from 'react-router-dom';
import { CircleItem, CirclePagination } from '../types';
import './style.scss'

interface Props {
  circles: CirclePagination
  onCircleClick: (id: string) => void
}

const CirclesSidebar = ({circles, onCircleClick}: Props) => {
  const [searchParams] = useSearchParams();

  return (
    <div className='circles__sidebar'>
      <div className='circles__sidebar__list'>
        {circles?.items.map(item => (
          <div className='circles__sidebar__list__item' onClick={() => onCircleClick(item.id)}>
            <FontAwesomeIcon icon={searchParams.has(item.id) ? "circle" : ["far", "circle"]} className="me-2" />
            {item.name}
          </div>
        ))}
      </div>
    </div>
  )
}

export default CirclesSidebar