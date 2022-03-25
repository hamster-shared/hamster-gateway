import type { AppRouteModule } from '/@/router/types';

import { LAYOUT } from '/@/router/constant';
import { t } from '/@/hooks/web/useI18n';

const dashboard: AppRouteModule = {
  path: '/gateway',
  name: 'Gateway',
  component: LAYOUT,
  redirect: '/gateway/initialization',
  meta: {
    orderNo: 10,
    icon: 'mdi:desktop-classic',
    title: t('routes.gateway.gateway'),
  },
  children: [
    {
      path: 'initialization',
      name: 'GatewayInitialization',
      component: () => import('/@/views/gateway/initialization/index.vue'),
      meta: {
        // affix: true,
        icon: 'mdi:cog-outline',
        title: t('routes.gateway.setting'),
      },
    },
    {
      path: 'boot',
      name: 'BootSetting',
      component: () => import('/@/views/gateway/boot/index.vue'),
      meta: {
        icon: 'bi:play-circle',
        title: t('routes.gateway.boot'),
      },
    },
  ],
};

export default dashboard;
