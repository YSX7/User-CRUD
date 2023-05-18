/* eslint-disable */

declare namespace NodeJS {
  interface ProcessEnv {
    AUTH_COOKIE: string;
    NODE_ENV: string;
    VUE_ROUTER_MODE: 'hash' | 'history' | 'abstract' | undefined;
    VUE_ROUTER_BASE: string | undefined;
  }
}
