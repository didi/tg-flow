
/**
 * Formatting time
 */
const formatDate = (value: any, fmt: any) => {
  fmt = fmt || 'YYYY-MM-DD HH:mm:ss'
  if (value === null) {
    return '-'
  } else {
    // return dayjs(formatISODate(value)).format(fmt)
  }
}
/**
 * Formatting iso date
 */
const formatISODate =( date: any) => {
  const [datetime, timezone] = date.split('+')
  if (!timezone || timezone.indexOf(':') >= 0) return date
  const hourOfTz = timezone.substring(0, 2) || '00'
  const secondOfTz = timezone.substring(2, 4) || '00'
  return `${datetime}+${hourOfTz}:${secondOfTz}`
}
/**
 * filter null
 */
const filterNull = (value: any) => {
  if (value === null || value === '') {
    return '-'
  } else {
    return value
  }
}

function dateFilter(date: any) {
  if (!date) {
    return '-'
  }

  const d = new Date(date)
  return d.toLocaleString()
}

export {
  formatDate, filterNull, dateFilter
}