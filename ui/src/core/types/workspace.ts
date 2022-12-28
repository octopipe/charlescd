export interface Workspace {
  name: string
  description: string
  routingStrategy: string
}

export interface WorkspaceModel extends Workspace {
  id: string
  createdAt: string
  circles: number
  modules: number
}

export interface WorkspacePagination {
  continue: string
  items: WorkspaceModel[]
}