import { AxiosResponse } from "axios";
import { Circle, CircleModel, CirclePagination, CircleStatusModel } from "../types/circle";
import { client } from "./client";

const getCircles = (workspaceId: string): Promise<AxiosResponse<CirclePagination>> => client.get(`/workspaces/${workspaceId}/circles`)
const getCircle = (workspaceId: string, circleId: string): Promise<AxiosResponse<CircleModel>> => client.get(`/workspaces/${workspaceId}/circles/${circleId}`)
const getCircleStatus = (workspaceId: string, circleId: string): Promise<AxiosResponse<CircleStatusModel>> => client.get(`/workspaces/${workspaceId}/circles/${circleId}/status`)
const createCircle = (workspaceId: string, data: Circle): Promise<AxiosResponse<CircleModel>> => client.post(`/workspaces/${workspaceId}/circles`, data)
const updateCircle = (workspaceId: string, circleId: string, data: Circle): Promise<AxiosResponse<CircleModel>> => client.put(`/workspaces/${workspaceId}/circles/${circleId}`, data)
const sync = (workspaceId: string, circleId: string): Promise<AxiosResponse<CircleStatusModel>> => client.post(`/workspaces/${workspaceId}/circles/${circleId}/sync`, { data: {} })
const getMetrics = (workspaceId: string, circleId: string): Promise<AxiosResponse<any>> => client.get(`/workspaces/${workspaceId}/circles/${circleId}/metrics`)
const metricQuery = (workspaceId: string, circleId: string, metricId: string, range: string): Promise<AxiosResponse<any>> => client.get(`/workspaces/${workspaceId}/circles/${circleId}/metrics/${metricId}/query?range=${range}`)


export const circleApi = {
  getCircles,
  getCircle,
  getCircleStatus,
  createCircle,
  updateCircle,
  sync,
  getMetrics,
  metricQuery,
}