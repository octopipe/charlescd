const jsonValidate = (value: string) => {
  try {
    JSON.parse(value)
  } catch (e: any) {
    return e?.message
  }
}


export const formValidations = {
  jsonValidate,
}