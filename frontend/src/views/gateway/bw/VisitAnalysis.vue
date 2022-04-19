<template>
  <div ref="chartRef" :style="{ height, width }"></div>
</template>
<script lang="ts">
  import { basicProps } from './props';
</script>
<script lang="ts" setup>
  import { onMounted,onUnmounted, ref, Ref } from 'vue';
  import { useECharts } from '/@/hooks/web/useECharts';
  import {getBandwidthApi} from '/@/api/gateway/bw'
  import { getBootStateApi } from '/@/api/gateway/boot';

  defineProps({
    ...basicProps,
  });
  const chartRef = ref<HTMLDivElement | null>(null);
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>);

  const xAxisLength = 18

  const rateIn = ref([0])
  const rateOut = ref([0])
  const timeArray = ref([new Date().toLocaleTimeString()])

  const timer = ref()


  const ensureData = function () {
    getBootStateApi().then(state => {

      if(state){
        getData()
      }

    })

  }

  const getData = function (){
    getBandwidthApi().then(data => {

      if(timeArray.value.length >= xAxisLength){
        timeArray.value.shift()
      }

      if(rateIn.value.length >= xAxisLength){
        rateIn.value.shift()
      }

      if(rateOut.value.length >= xAxisLength){
        rateOut.value.shift()
      }

      timeArray.value.push(new Date().toLocaleTimeString())
      rateIn.value.push(data.RateIn)
      rateOut.value.push(data.RateOut)


      setOptions({
        animation: false,
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            lineStyle: {
              width: 1,
              color: '#019680',
            },
          },
        },
        xAxis: {
          type: 'category',
          boundaryGap: false,
          data: timeArray,
          splitLine: {
            show: true,
            lineStyle: {
              width: 1,
              type: 'solid',
              color: 'rgba(226,226,226,0.5)',
            },
          },
          axisTick: {
            show: false,
          },
        },
        yAxis: [
          {
            type: 'value',
            max: 80000,
            splitNumber: 4,
            axisTick: {
              show: false,
            },
            splitArea: {
              show: true,
              areaStyle: {
                color: ['rgba(255,255,255,0.2)', 'rgba(226,226,226,0.2)'],
              },
            },
          },
        ],
        grid: { left: '1%', right: '1%', top: '2  %', bottom: 0, containLabel: true },
        series: [
          {
            smooth: true,
            data: rateIn,
            type: 'line',
            areaStyle: {},
            itemStyle: {
              color: '#5ab1ef',
            },
          },
          {
            smooth: true,
            data: rateOut,
            type: 'line',
            areaStyle: {},
            itemStyle: {
              color: '#019680',
            },
          },
        ],
      });

    })
  }

  onMounted(() => {

    timer.value = setInterval(ensureData,2000);


  });

  onUnmounted(() => {
    clearInterval(timer.value)
  });
</script>
