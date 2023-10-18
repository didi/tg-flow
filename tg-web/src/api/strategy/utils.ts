interface VauleAndLabelI {
  value: string;
  label: string;
}

import {
  scenesConfigAppnameContentI,
  scenesConfigNameContentI,
  experimentHistoryVersionResFilterI,
  experimentHistoryVersionResI,
} from "./types";
export function getTime(params: number): number {
  let date = new Date(params * 1000);
  return Number(
    `${date.getFullYear()}${
      date.getMonth() + 1 < 10
        ? "0" + (date.getMonth() + 1)
        : date.getMonth() + 1
    }${date.getDate() < 10 ? "0" + date.getDate() : date.getDate()}${
      date.getHours() < 10 ? "0" + date.getHours() : date.getHours()
    }${date.getMinutes() < 10 ? "0" + date.getMinutes() : date.getMinutes()}${
      date.getSeconds() < 10 ? "0" + date.getSeconds() : date.getSeconds()
    }`
  );
}

export function scenesSelectAppName(
  params: Array<scenesConfigAppnameContentI>
): Array<VauleAndLabelI> {
  const arr: Array<VauleAndLabelI> = [];  
  params.forEach((p) => {
    arr.push({
      label: p.AppName,
      value: p.AppName,
    });
  });
  arr.shift();
  return arr;
}

export function scenesSelectName(
  params: Array<scenesConfigNameContentI>
): Array<VauleAndLabelI> {
  const arr: Array<VauleAndLabelI> = [];
  params.forEach((p) => {
    arr.push({
      label: p.SceneName,
      value: p.SceneName,
    });
  });
  return arr;
}

export function experimentAppname(
  params: Array<scenesConfigAppnameContentI>
): Array<VauleAndLabelI> {
  const arr: Array<VauleAndLabelI> = [];
  params.forEach((p) => {
    arr.push({
      label: p.AppName,
      value: p.AppId,
    });
  });
  return arr;
}

export function filterData(
  params: Array<experimentHistoryVersionResI>
) :Array<experimentHistoryVersionResFilterI> {
  return params.map((p) => {
    return {
      dimension_id:p.dimension_id,
      version_create_time:p.version_create_time,
      version_id:p.version_id
    }
  })
}

export function algoName(params: any): Array<VauleAndLabelI>{
  const arr: Array<VauleAndLabelI> = [];
  params.forEach((p : any) => {
    if(p.AlgoName){
      arr.push({
        label:p.AlgoName,
        value:p.AlgoName
      })
    }
  });
  return arr
}

export function modelName(params:any[]): Array<VauleAndLabelI>{
  const arr: Array<VauleAndLabelI> = [];
  params.forEach((p : any) => {
    if(p.ModelName){
      arr.push({
        label:p.ModelName,
        value:p.ModelName
      })
    }
  });
  return arr
}
export function appName(params:Object[]): Array<VauleAndLabelI>{
  const arr: Array<VauleAndLabelI> = [];  
  params.forEach((p : any) => {
    if(p.AppName){
      arr.push({
        label:p.AppName,
        value:p.AppId
      })
    }
  });
  return arr
}
