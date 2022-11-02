export interface TemplateForm {
  label: string
  jsonKey: string
}

export interface Template {
  name: string
  description: string
  icon: string
  form: TemplateForm[]
}

export const TEMPLATES: Template[] = [
  {
    name: 'Angular',
    description: 'Template for angular projects',
    icon: 'fa-brands fa-angular',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'React',
    description: 'Template for react projects',
    icon: 'fa-brands fa-react',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Springboot',
    description: 'Template for Springboot projects',
    icon: 'fa-brands fa-java',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Micronaut',
    description: 'Template for Micronaut projects',
    icon: 'fa-brands fa-java',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Nestjs',
    description: 'Template for Node projects',
    icon: 'fa-brands fa-node',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Node app',
    description: 'Template for Node projects',
    icon: 'fa-brands fa-node',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Golang app',
    description: 'Template for golang app',
    icon: 'fa-brands fa-golang',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Static app',
    description: 'Template for static app',
    icon: 'fa-brands fa-html5',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  },
  {
    name: 'Custom app',
    description: 'Template for custom app',
    icon: 'fa-brands fa-docker',
    form: [
      {
        label: 'Container image',
        jsonKey: '$.spec.template.spec.containers[0].image'
      }
    ]
  }
]