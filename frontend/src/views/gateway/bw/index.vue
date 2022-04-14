<template>
  <dev>
    <PageWrapper title="bandwidth">
      <template #headerContent>
        <div class="flex justify-between items-center">
        <span class="flex-1">
          <a href="#" target="_blank">{{ name }}</a>
          {{ t('boot.boot.bootTitle') }}
        </span>
        </div>
      </template>

      <VisitAnalysis />
    </PageWrapper>
  </dev>
</template>

<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import { PageWrapper } from '/@/components/Page';
import VisitAnalysis from './VisitAnalysis.vue'
import { useI18n } from '/@/hooks/web/useI18n';
const { t } = useI18n();

const { pkg } = __APP_INFO__;
const { name } = pkg;
</script>

<style scoped>

</style>
