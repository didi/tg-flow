export function deepcopy(o :any) {
  return JSON.parse(JSON.stringify(o))
}

export function copyToClipboard(value: any, successfully: any, failure: any) {
  const clipboard = navigator.clipboard
  if (clipboard !== undefined) {
    navigator.clipboard.writeText(value).then(successfully, failure)
  } else {
    // fallback to execCommand
    let isSuccess = false

    if (document.queryCommandSupported && document.queryCommandSupported('copy')) {
      const el = document.createElement('textarea')
      el.value = value
      el.style.top = '0'
      el.style.left = '0'
      el.style.position = 'fixed' // Prevent scrolling to bottom of page in Microsoft Edge.
      document.body.appendChild(el)

      el.focus()
      el.select()

      // Security exception may be thrown by some browsers.
      try {
        if (document.execCommand('copy')) {
          isSuccess = true
        }
      } catch (ex) {
        console.warn('Copy to clipboard failed.', ex)
      } finally {
        document.body.removeChild(el)
      }
    }

    // callback
    if (isSuccess) {
      if (successfully !== undefined) {
        successfully()
      }
    } else {
      if (failure !== undefined) {
        failure()
      }
    }
  }
}