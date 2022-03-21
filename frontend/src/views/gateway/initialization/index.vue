<template>
  <div>
    <PageWrapper :title="t('initialization.initialization.basicConfiguration')">
      <CollapseContainer :title="t('initialization.initialization.chainSettings')" :canExpan="false">
        <a-row :gutter="24">
          <a-col :span="18">
            <BasicForm @register="chainRegister" />
          </a-col>
        </a-row>
      </CollapseContainer>

      <Button class="mt-4" type="primary" @click="handleSubmit"> {{ t('initialization.initialization.updateInformation') }} </Button>
    </PageWrapper>
    <a-modal
      v-model:visible="visible"
      :title="t('initialization.initialization.addGatewayNode')"
      :maskClosable="false"
      :footer="null"
      :centered="true"
      :closable="false"
    >
      <a-spin :spinning="addLoading">
        <div class="staking-content">
          <span class="title">{{ t('initialization.initialization.gatewayNode') }}</span>
          <a-textarea
            v-model:value="value"
            :placeholder="t('initialization.initialization.inputGatewayNodeTip')"
            :rows="3"
            @change="checkAddBootstrap"
            style="width: 420px"
          />
        </div>
        <span class="form-error-tip" v-if="addBootstrapTip">{{
            t('initialization.initialization.gatewayNodeTip')
          }}</span>
        <div class="staking-footer">
          <a-button class="staking-btn-close" @click="close">{{
              t('accountInfo.info.cancel')
            }}</a-button>
          <a-button class="staking-btn-ok" @click="ok">{{
              t('accountInfo.info.determine')
            }}</a-button>
        </div>
      </a-spin>
    </a-modal>
  </div>
</template>
<script lang="ts">
import {computed, defineComponent, reactive, onMounted, h, ref, toRefs} from 'vue';
  import { Button, Row, Col } from 'ant-design-vue';
  import { BasicForm, useForm } from '/@/components/Form/index';
  import { CollapseContainer } from '/@/components/Container';
  import { PageWrapper } from '/@/components/Page';
  import { Description} from '/@/components/Description';
  import { useI18n } from '/@/hooks/web/useI18n';
  import { useMessage } from '/@/hooks/web/useMessage';

  import headerImg from '/@/assets/images/header.jpg';
  import { getConfigApi,setConfigApi } from '/@/api/gateway/initialization';
  import { chainSchemas } from './data';
  import { useUserStore } from '/@/store/modules/user';
  import {ProviderConfig} from "/@/api/gateway/model/settingModel";
  import AButton from "/@/components/Button/src/BasicButton.vue";

  export default defineComponent({
    components: {
      AButton,
      BasicForm,
      CollapseContainer,
      Button,
      ARow: Row,
      ACol: Col,
      PageWrapper,
      Description,
    },
    setup: function () {
      const {createMessage} = useMessage();
      const userStore = useUserStore();
      const {t} = useI18n();
      const state = reactive({
        visible: false,
        addLoading: false,
        addBootstrapTip: false,
        value: "",
      });

      const [chainRegister, {
        setFieldsValue: chainSetFieldsValue,
        validateFields: chainValidateFields
      }] = useForm({
        labelWidth: 120,
        schemas: chainSchemas,
        showActionButtonGroup: false,
      })


      onMounted(async () => {
        const data = await getConfigApi();
        await chainSetFieldsValue({
          "address": data.chainApi,
          "account": data.seedOrPhrase,
          "publicIp": data.publicIp,
          "publicPort": data.publicPort,
        });
      });

      const avatar = computed(() => {
        const {avatar} = userStore.getUserInfo;
        return avatar || headerImg;
      });

      return {
        t,
        avatar,
        chainRegister,
        ...toRefs(state),
        handleSubmit: () => {
          Promise.all([chainValidateFields()])
            .then((data) => {
              let chainValues = data[0];
              let config: ProviderConfig = {
                chainApi: chainValues.address,
                seedOrPhrase: chainValues.account,
                publicIp: chainValues.publicIp,
                publicPort: chainValues.publicPort,

              };
              setConfigApi(config)
                .then(() => {
                  createMessage.success(t('initialization.initialization.updateSucceeded'));
                })
                .catch((err) => {
                  createMessage.error(t('initialization.initialization.updateFailed'), err);
                });
            })
            .catch((err) => {
              createMessage.error(t('initialization.initialization.verificationFailed'), err);
            });
        },
      };
    },
  });
</script>

<style lang="less" scoped>
  .strap-style {
    display: flex;
    align-items: center;
  }
  .staking-content {
    display: flex;
    align-items: center;
    margin-top: 24px;
    padding: 0px 16px;
    .title {
      min-width: 40px;
      margin-right: 8px;
      color: rgba(0, 0, 0, 0.85);
    }
  }
  .form-error-tip {
    color: #f5313d;
    font-style: normal;
    font-weight: normal;
    font-size: 10px;
    line-height: 17px;
    margin-left: 80px;
  }
  .staking-footer {
    margin-top: 24px;
    display: grid;
    padding: 0px 16px 24px 16px;
    grid-template-columns: 1fr 1fr;
    gap: 12px;
    .staking-btn-close {
      width: 100%;
    }
    .staking-btn-ok {
      background-color: rgb(24, 144, 255);
      color: white;
    }
  }
</style>
