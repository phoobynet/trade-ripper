import axios, { AxiosRequestConfig } from 'axios'

let baseUrl = '/api'

if (import.meta.env.DEV) {
  baseUrl = `http://${location.hostname}:3000/api`
}

export const http = axios.create({
  baseURL: baseUrl,
})

export const getData = <T>(
  url: string,
  config: AxiosRequestConfig = {},
): Promise<T> => {
  return http.get<T>(url, config).then((response) => response.data)
}
