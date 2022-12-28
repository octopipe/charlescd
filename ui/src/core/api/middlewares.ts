export const handleErrorsMiddleware = ({ dispatch, getState }: any) => (next: any) => (action: any) => {

  return next(action);
}