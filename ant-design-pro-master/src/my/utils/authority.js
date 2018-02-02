const TOKEN = 'token';

export function getAuthority() {
  return localStorage.getItem(TOKEN);
}

export function setAuthority(authority) {
  return localStorage.setItem(token, authority);
}
