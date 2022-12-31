import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import React, { useEffect, useState } from "react";
import { Nav } from "react-bootstrap";
import { useParams } from "react-router-dom";
import { Area, AreaChart, CartesianGrid, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { circleApi } from "../../../core/api/circle";
import ViewInput from "../../../core/components/ViewInput";
import { CircleModel } from "../../../core/types/circle";
import './style.scss'

interface Props {
  circle?: CircleModel 
}

enum BASIC_METRIC_TYPE {
  CPU = 'CPU',
  MEMORY = 'MEMORY'
}

const CircleMetrics = ({ circle }: Props) => {
  const { workspaceId } = useParams()
  const [basicMetricType, setBasicMetricType] = useState<string>('')
  const [metrics, setMetrics] = useState<any>([])
  const [metricResult, setMetricResult] = useState([])

  const getMetrics = async () => {
    const res = await circleApi.getMetrics(workspaceId || '', circle?.id || '')
    setMetrics(res?.data?.items)
  }

  const metricQuery = async (metricId: string) => {
    const res = await circleApi.metricQuery(workspaceId || '', circle?.id || '', metricId)
    if (res?.data?.length > 0) {
      const values = res?.data[0].values || []
      setMetricResult(values.map((v: any) => {
        const date = new Date(v[0] *1000)
        const hour = '0' + date.getHours()
        const minutes = '0' + date.getMinutes()
        const seconds = date.getSeconds()

        return {time: `${hour.slice(-2)}:${minutes.slice(-2)}:${seconds}`, usage: v[1]}
      }))
    }
  }

  useEffect(() => {
    getMetrics()
  }, [])

  useEffect(() => {
    if (metrics.length <= 0)
      return

    setBasicMetricType(metrics[0].id)
  }, [metrics])

  useEffect(() => {
    metricQuery(basicMetricType)
    const interval = setInterval(() => {
      metricQuery(basicMetricType)
    }, 3000)
    
    return () => clearInterval(interval)
  }, [basicMetricType])

  return (
    <div className="circle-metrics">

      <ViewInput
        label="Basic"
        icon="chart-simple"
      >
        <Nav variant='pills' defaultActiveKey={basicMetricType} onSelect={key => setBasicMetricType(key || '')} className='mb-3'>
          {metrics?.map((metric: any) => (
            <Nav.Item><Nav.Link eventKey={metric.id}>{metric.name}</Nav.Link></Nav.Item>
          ))}
        </Nav>
          {metrics?.filter((metric: any) => metric.id == basicMetricType)?.map((metric: any) => (
            <div className="circle-metrics__base">
              <ResponsiveContainer width="100%" height="100%">
                <AreaChart
                  width={500}
                  height={400}
                  data={metricResult}
                  margin={{
                    top: 10,
                    right: 30,
                    left: 0,
                    bottom: 0,
                  }}
                >
                  <XAxis dataKey="time" />
                  <YAxis />
                  <Tooltip />
                  <Area type="monotone" dataKey="usage" stackId="1" stroke="#8884d8" fill="#8884d8" />
                </AreaChart>
              </ResponsiveContainer>
            </div>
          ))}
      </ViewInput>
    </div>
  )
}

export default CircleMetrics