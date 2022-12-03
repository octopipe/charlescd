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