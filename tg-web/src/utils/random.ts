export function Random(min: number, max: number) {
  return Math.round(Math.random() * (max - min)) + min
}

export function getRandomCode() {
  return Random(0, 9999999)
}
