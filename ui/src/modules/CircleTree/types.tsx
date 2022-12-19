export interface ResourceOwner {
  name: string
  kind: string
  controller?: boolean
}

export interface ResourceMetadata {
  name: string
  namespace: string
  kind: string
  group: string
  owner?: ResourceOwner
  status?: string
  message?: string
}

export interface Resource {
  metadata: ResourceMetadata
  manifest: any
}