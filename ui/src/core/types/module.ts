export interface Module {
  name: string
  description: string
  path: string
  url: string
  templateType: string
  visibility: string
}

export interface ModuleModel extends Module {
  id: string
  createdAt: string
}

export interface ModulePagination {
  continue: string
  items: ModuleModel[]
}