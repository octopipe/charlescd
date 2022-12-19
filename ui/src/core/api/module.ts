import { AxiosResponse } from "axios";
import { Circle, CircleModel, CirclePagination } from "../types/circle";
import { ModulePagination } from "../types/module";
import { client } from "./client";

const getModules = (workspaceId: string): Promise<AxiosResponse<ModulePagination>> => client.get(`/workspaces/${workspaceId}/modules`)

export const moduleApi = {
  getModules,
}