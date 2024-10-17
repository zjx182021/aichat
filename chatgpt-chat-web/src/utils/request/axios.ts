import axios, { type AxiosResponse } from 'axios'
import { deleteCookieByKey, getCookieValue } from '../cookie/index'

const service = axios.create({
  baseURL: import.meta.env.VITE_GLOB_API_URL,
})

service.interceptors.request.use(
  (config) => {
    const access_token = getCookieValue('sso_0voice_access_token')
    if (access_token)
      config.headers.Authorization = access_token
    return config
    // else
    // window.location.href = import.meta.env.VITE_USER_CENTER
    // return config
  },
  (error) => {
    return Promise.reject(error.response)
  },
)

service.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse => {
    if (response.status === 200)
      return response

    if (response.status === 401) {
      deleteCookieByKey('sso_0voice_access_token')
      window.location.href = import.meta.env.VITE_USER_CENTER
    }

    throw new Error(response.status.toString())
  },
  (error) => {
    return Promise.reject(error)
  },
)

export default service
