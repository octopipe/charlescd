import { CircleItemModule } from "../CircleModules/types"


export interface CircleItem {
  id: string
  name: string
  namespace: string
  isDefault: boolean
  modules: CircleItemModule[]
}

export interface CirclePagination {
  continue: string
  items: CircleItem[]
}