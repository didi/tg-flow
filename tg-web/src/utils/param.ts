
const REF_KEY = '$ref'

export function fromRef(o: any) {
  if (o[REF_KEY] === undefined) {
    o[REF_KEY] = ''
  }
  return o[REF_KEY]
}

export function toRef(s: any) {
  return {
    [REF_KEY]: s
  }
}
