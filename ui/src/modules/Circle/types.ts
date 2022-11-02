import { CircleItem } from "../Circles/types";

export interface CircleRouting {
  strategy: string
}

export interface Overrides {
  [key: string]: string
}

export interface ModuleStatusValue {
  status: string
  error: string
}

export interface ModuleStatus {
  [key: string]: ModuleStatusValue
}

export interface CircleStatus {
  modules: ModuleStatus
}
 
export interface CircleModule {
  name: string
  revision: string
  overrides: Overrides
}

export interface CircleEnrivonment {
  key: string
  value: string
}

export interface Circle {
  metadata: CircleItem
  routing: CircleRouting
  modules: CircleModule[]
  status: CircleStatus
  environments: CircleEnrivonment[]
}