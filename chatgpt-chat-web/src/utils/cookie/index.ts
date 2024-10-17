export function getCookieValue(key: string) {
  const cookies = document.cookie.split(';')
  for (let i = 0; i < cookies.length; i++) {
    const cookie = cookies[i].trim()
    if (cookie.startsWith(`${key}=`))
      return cookie.substring(key.length + 1)
  }
  return null
}

export function deleteCookieByKey(key: string) {
  document.cookie = `${key}=;expires=Thu, 01 Jan 1970 00:00:00 GMT;path=/`
}
