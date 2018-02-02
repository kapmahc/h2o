import RenderAuthorized from 'ant-design-pro/lib/Authorized';

const {Secured} = RenderAuthorized('user');
const TOKEN = "token"

export const getToken = () => {
  return localStorage.getItem(TOKEN)
}
export const setToken = (t) => {
  if (t) {
    localStorage.setItem(TOKEN, t)
  } else {
    localStorage.removeItem(TOKEN);
  }
}

// export Secured;