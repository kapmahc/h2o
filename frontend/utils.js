export const backend = (u) => `${process.env.BACKEND}${u}`;

const LOCALE = 'locale';
export const get_locale = () => window.localStorage.getItem(LOCALE) || 'en-US';
export const set_locale = (l) => window.localStorage.setItem(LOCALE, l);
