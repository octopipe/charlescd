export interface CircleRoutingCustomMatchHeaders {
  [key: string]: string
}

export interface CircleRoutingCustomMatch {
  headers: CircleRoutingCustomMatchHeaders
}

export interface CircleRoutingSegment {
  key: string
  op: string
  value: string
}

export interface CircleRoutingMatch {
  customMatch?: CircleRoutingCustomMatch
  segments?: CircleRoutingSegment[]
}

export interface CircleRoutingCanary {
  weight: number
}

export interface CircleRouting {
  strategy: string
  match?: CircleRoutingMatch
  canary?: CircleRoutingCanary
}

export interface Overrides {
  [key: string]: string
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

export interface CircleStatus {
  modules: ModuleStatus
}
 
export interface CircleModule {
  moduleId: string
  name: string
  revision: string
  overrides: Overrides
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

export enum CIRCLE_VIEW_MODE {
  VIEW = 'VIEW',
  CREATE = 'CREATE',
  UPDATE = 'UPDATE',
  TREE = 'TREE',
}

export enum CIRCLE_MODULE_RESOURCE_STATUS {
  HEALTHY = 'Healthy',
  PROGRESSING = 'Progressing',
  DEGRADED = 'Degraded',
}