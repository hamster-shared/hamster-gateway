import { FormSchema } from '/@/components/Form';
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();
// 基础设置 form
export const chainSchemas: FormSchema[] = [
  {
    field: 'address',
    component: 'Select',
    label: 'chainAddress',
    colProps: { span: 12 },
    componentProps: {
      options: [
        { label: 'dev(ws://127.0.0.1:9944)', value: 'ws://127.0.0.1:9944' },
        { label: 'test(ws://183.66.65.207:49944)', value: 'ws://183.66.65.207:49944' },
      ],
    },
    rules: [{ required: true }],
  },
  {
    field: 'account',
    component: 'Input',
    label: 'account',
    colProps: { span: 12 },
    componentProps: {
      placeholder: t('initialization.initialization.seedTip'),
    },
    rules: [{ required: true }],
  },
  {
    field: 'publicIp',
    component: 'Input',
    label: 'publicIp',
    colProps: {span: 12},
    componentProps: {
      placeholder: t('initialization.initialization.publicIp'),
    },
    rules: [{required: true}],
  },
  {
    field: 'publicPort',
    component: 'InputNumber',
    label: 'publicPort',
    colProps: {span: 12},
    componentProps: {
      placeholder: t('initialization.initialization.publicPort'),
    }
  },
]
