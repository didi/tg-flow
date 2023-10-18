export function debounce(fn : Function, delay : number)  {
  let timer = null as any;
  return function(...arg : any){
    if(timer) clearTimeout(timer)
    timer && clearTimeout(timer)
    timer = setTimeout(() => {
      //@ts-ignore
      fn.call(this,arg)
    },delay)
  }
}