export interface CircleRoutingMatchHeaders {
  [key: string]: string
}

export interface CircleRoutingMatch {
  headers: CircleRoutingMatchHeaders
}

export interface CircleRoutingSegment {
  key: string
  op: string
  value: string
}

export interface CircleRoutingCanary {
  weight: number
}

export interface CircleRouting {
  strategy: string
  match?: CircleRoutingMatch
  canary?: CircleRoutingCanary
  segments?: CircleRoutingSegment[]
}

export interface Override {
  key: string
  value: string
}

export interface ModuleResource {
  group: string
  kind: string
  name: string
  namespace: string
  status?: string
  error?: string
}

export interface ModuleStatusValue {
  status: string
  error: string
  resources: ModuleResource[]
}

export interface ModuleStatus {
  [key: string]: ModuleStatusValue
}

export interface CircleStatusHistory {
  action: string
  eventTime: string
  message: string
  status: string
}

export interface CircleStatus {
  modules: ModuleStatus
  history: CircleStatusHistory[]
}
 
export interface CircleModule {
  name: string
  revision: string
  overrides: Override[]
}

export interface CircleEnrivonment {
  key: string
  value: string
}

export interface Circle {
  name: string
  description: string
  routing: CircleRouting
  modules: CircleModule[]
  environments: CircleEnrivonment[]
}

export interface CircleModel extends Circle {
  id: string
  createdAt: string
  status: CircleStatus
}

export interface CircleItem {
  id: string
  name: string
  namespace: string
  isDefault: boolean
  modules: CircleModule[]
  routing: CircleRouting
}

export interface CircleStatusModelModuleResource {
  name: string
  namespace: string
  kind: string
  health?: string
  message?: string
}

export interface CircleStatusModelModule {
  resources: CircleStatusModelModuleResource[]
} 

export interface CircleStatusModelModules {
  [key: string]: CircleStatusModelModule
}

export interface CircleStatusModel {
  modules: CircleStatusModelModules
}

export interface CirclePagination {
  continue: string
  items: CircleItem[]
}

export enum CIRCLE_ROUTING_STRATEGY {
  MATCH = 'MATCH',
  SEGMENTS = 'SEGMENTS',
  CANARY = 'CANARY'
}

export enum CIRCLE_VIEW_MODE {
  VIEW = 'VIEW',
  CREATE = 'CREATE',
  UPDATE = 'UPDATE',
  TREE = 'TREE',
  HISTORY = 'HISTORY',
  METRICS = 'METRICS'
}

export enum CIRCLE_MODULE_RESOURCE_STATUS {
  HEALTHY = 'Healthy',
  PROGRESSING = 'Progressing',
  DEGRADED = 'Degraded',
}