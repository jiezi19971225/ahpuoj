<template lang="pug">
.admin-content
  .content__main
    .chart__wrapper
      h2.pl20.pt10 近期提交情况
      line-chart(:option="chartOption",:id="'chart'",style="width:100%;height:500px;")
</template>

<script>
import LineChart from 'common/components/LineChart/index.vue';
import { getSubmitStatistic } from 'admin/api/admin';

export default {
  components: {
    LineChart,
  },
  data() {
    return {
      chartOption: {
        color: ['#ffdf25'],
        // title: {
        //   text: "123"
        // },
        tooltip: {},
        legend: {
          data: ['累计提交'],
        },
        xAxis: {
          type: 'time',
        },
        yAxis: {},
        series: [
          {
            name: '累计提交',
            type: 'line',
            data: [],
          },
        ],
      },
    };
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      try {
        const { id } = this.$route.params;
        const res = await getSubmitStatistic(id);
        const { data } = res;
        this.chartOption.series[0].data = data.recent_submit_statistic;
      } catch (err) {
        console.log(err);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
.chart__wrapper {
  background: $--color-level15;
}
</style>
