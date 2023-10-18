import axios from "axios";
import qs from "query-string";

// axios.defaults.baseURL = "/api";
// axios的封装
export async function request<T = any>(
  path: string,
  params: object,
  method: string,
  type?: number
) {
  let content = "x-www-form-urlencoded";
  let data = qs.stringify(params);
  if (type === 1) {
    content = "json";
    data = JSON.stringify(params);
  }
  try {
    const serve = await axios({
      method,
      url: path,
      withCredentials: true,
      headers: {
        "Content-Type": `application/${content}`,
      },
      data,
    });
    if (serve.data) {
      return serve.data as T;
    } else {
      return null;
    }
  } catch (err: any) {
    if (err.response.status === 302 || err.response.status === 301) {
      // TODO 登录认证
      window.location.replace("/main");
    }
    return null;
  }
}

export default request;
