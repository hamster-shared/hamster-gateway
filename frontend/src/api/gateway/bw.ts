import { defHttp } from '/@/utils/http/axios';
import {BandwidthStats} from "/@/api/gateway/model/bwModel";

enum Api {
  BW = '/v1/p2p/bw',
}

// get bandwidth
// get boot state
export const getBandwidthApi = () => {
  return defHttp<BandwidthStats>.get({url:Api.BW})
}
