import { AxiosResponse } from "axios";
import { Workspace, WorkspaceModel, WorkspacePagination } from "../types/workspace";
import { client } from "./client";

const getWorkspaces = (): Promise<AxiosResponse<WorkspaceModel[]>> => client.get(`/workspaces`)
const createWorkspace = (data: Workspace): Promise<AxiosResponse<WorkspaceModel>> => client.post(`/workspaces`, data)
const updateWorkspace = (workspaceId: string, data: Workspace): Promise<AxiosResponse<WorkspaceModel>> => client.put(`/workspaces/${workspaceId}`, data)


export const workspaceApi = {
  getWorkspaces,
  createWorkspace,
  updateWorkspace,
}