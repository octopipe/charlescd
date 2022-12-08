export interface Workspace {
  name: string
  description: string
  routingStrategy: string
}

export interface WorkspaceModel extends Workspace {
  id: string
  createdAt: string
}