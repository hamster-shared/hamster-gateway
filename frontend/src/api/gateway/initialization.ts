import { defHttp } from '/@/utils/http/axios';
import {GatewayConfig,} from '/@/api/gateway/model/settingModel';

enum Api {
  Setting = '/api/v1/config/settting',
}

//获取系统配置
export const getConfigApi =  () => {
  return defHttp.get<GatewayConfig>({ url: Api.Setting });
}

// 修改配置
export const setConfigApi = (config: GatewayConfig) => {
  return defHttp.post({ url: Api.Setting, data: config })
}
