import React, { useEffect } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import { useSearchParams } from 'react-router-dom';
import { CirclePagination } from '../types';
import './style.scss'
import { Nav } from 'react-bootstrap';
import Spinner from '../../../core/components/Spinner';
import DynamicContainer from '../../../core/components/DynamicContainer';

interface Props {
  circles?: CirclePagination
  loading: boolean
  onCircleClick: (id: string) => void
  onCircleCreateClick: () => void
}

const CirclesSidebar = ({circles, loading, onCircleClick, onCircleCreateClick}: Props) => {
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
      <DynamicContainer loading={loading} className='circles__sidebar__list'>
        {circles && circles.items.length > 0 && circles?.items.map(item => (
          <div key={item.id} className='circles__sidebar__list__item' onClick={() => onCircleClick(item.id)}>
            <FontAwesomeIcon icon={searchParams.has(item.id) ? "circle" : ["far", "circle"]} className="me-2" />
            {item.name}
          </div>
        ))}
        {circles && circles.items.length <= 0 && (
          <div className='text-muted'>
            There are no circles here
          </div>
        )}
      </DynamicContainer>
    </div>
  )
}

export default CirclesSidebar