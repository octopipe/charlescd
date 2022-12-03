import { IconProp } from '@fortawesome/fontawesome-svg-core';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import React, { useState } from 'react';
import { Link, NavLink, useMatch, useParams, useSearchParams } from 'react-router-dom';
import { CircleItem, CirclePagination } from '../types';
import { ReactComponent as EmptyCirclesSVG } from '../../../core/assets/svg/empty-circles.svg'
import './style.scss'
import Placeholder from '../../../core/components/Placeholder';
import { Nav } from 'react-bootstrap';

interface Props {
  circles: CirclePagination
  onCircleClick: (id: string) => void
  onCircleCreateClick: () => void
}

const CirclesSidebar = ({circles, onCircleClick, onCircleCreateClick}: Props) => {
  const [searchParams] = useSearchParams();

  return (
    <div className='circles__sidebar'>
      <div className='circles__sidebar__header'>
        <Nav>
          <Nav.Item onClick={onCircleCreateClick}>
            <FontAwesomeIcon icon="plus-circle" className="me-1" /> Create circle
          </Nav.Item>
        </Nav>
      </div>
      <div className='circles__sidebar__list'>
        {circles.items.length > 0 ? (
          <>
            {circles?.items.map(item => (
            <div key={item.id} className='circles__sidebar__list__item' onClick={() => onCircleClick(item.id)}>
              <FontAwesomeIcon icon={searchParams.has(item.id) ? "circle" : ["far", "circle"]} className="me-2" />
              {item.name}
            </div>
          ))}
          </>
        ) : (
          <div className='text-muted'>
            There are no circles here
          </div>
        )}
      </div>
    </div>
  )
}

export default CirclesSidebar