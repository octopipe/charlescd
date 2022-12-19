import { AxiosResponse } from "axios";
import { Circle, CircleModel, CirclePagination } from "../types/circle";
import { client } from "./client";

const getCircles = (workspaceId: string): Promise<AxiosResponse<CirclePagination>> => client.get(`/workspaces/${workspaceId}/circles`)
const getCircle = (workspaceId: string, circleId: string): Promise<AxiosResponse<CircleModel>> => client.get(`/workspaces/${workspaceId}/circles/${circleId}`)

export const circleApi = {
  getCircles,
  getCircle,
}