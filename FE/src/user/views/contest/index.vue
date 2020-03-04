<template lang="pug">
  .content
    title {{`${contest?contest.name:''} - AHPUOJ`}}
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            .header 竞赛信息
          li  状态：
              template(v-if="contest")
                el-tag(v-if="contest.status==1",type="success",effect="dark") 未开始
                el-tag(v-if="contest.status==2",type="primary",effect="dark") 进行中
                el-tag(v-if="contest.status==3",type="danger",effect="dark") 已结束
          li 模式：
              template(v-if="contest")
                el-tag(:type="contest.private == 1 ? 'danger':'success'",effect="dark")  {{ contest.private == 1?"私有赛":"公开赛" }}
                el-tag(:type="contest.team_mode == 0 ? 'success':'primary'",effect="dark")  {{ contest.team_mode == 0?"个人赛":"团体赛" }}
          li
            div 开始时间：
            p(v-if="contest") {{contest.start_time}}
            div.mt10 结束时间：
            p(v-if="contest")  {{contest.end_time}}
          li  持续时长：
            p(v-if="contest")  {{timeDiff}}
        .button__wrapper
          el-button(size="small",type="primary",@click="jumpToStatus") 记录
          el-button(size="small",type="primary",@click="jumpToRank") 排名
          el-button(v-if="contest&&contest.team_mode == 1",size="small",type="primary",@click="jumpToTeamRank") 团队排名
      .main
        h1 {{contest?contest.name:''}}
        .main__section(style="min-height:200px;")
          h3 竞赛简介
          div(v-if="contest",v-html="contest.description")
        .main__section
          h3 问题列表
          el-table.dataTable(v-if="seeable",:data="contest.probleminfos", style="width: 100%")
            el-table-column(width="40")
              template(slot-scope="scope")
                svg-icon(name="ok",v-if="scope.row.status == 1")
                svg-icon(name="wrong",v-else-if="scope.row.status == 2")
            el-table-column(label="#", width="60")
              template(slot-scope="scope")
                span {{ engNum(scope.row.num) }}
            el-table-column(label="标题", min-width="180")
              template(slot-scope="scope")
                router-link(:to="{name:'contestProblem',params:{id:contest.id,num:scope.row.num}}", target="_blank") {{scope.row.title}}
          p(v-else) {{reason}}
</template>

<script>
import { getContest, getLanguageList } from 'user/api/nologin';

import EventBus from 'common/eventbus';
import { submitJudgeCode } from 'user/api/user';

export default {
  data() {
    return {
      seeable: false,
      reason: '',
      contest: null,
      langList: [],
    };
  },
  computed: {
    contestStatus() {
      const startDate = new Date(this.contest.start_time);
      const endDate = new Date(this.contest.end_time);
      const nowDate = new Date();
      if (nowDate.getTime() < startDate.getTime()) {
        // 未开始
        return 0;
      } if (nowDate.getTime() > endDate.getTime()) {
        // 已结束
        return 2;
      }
      // 进行中
      return 1;
    },
    timeDiff() {
      const startDate = new Date(this.contest.start_time);
      const endDate = new Date(this.contest.end_time);
      const dateDiff = endDate.getTime() - startDate.getTime();

      const days = Math.floor(dateDiff / (24 * 3600 * 1000));
      let left = dateDiff % (24 * 3600 * 1000);
      const hours = Math.floor(left / (3600 * 1000));
      left %= (3600 * 1000);
      const minutes = Math.floor(left / (60 * 1000));
      left %= (60 * 1000);
      const seconds = Math.round(left / 1000);

      let res = '';
      res += days ? `${days}天` : '';
      res += hours ? `${hours}小时` : '';
      res += minutes ? `${minutes}分钟` : '';
      res += seconds ? `${seconds}秒` : '';
      return res;
    },
  },
  mounted() {
    this.init();
  },
  methods: {
    async init() {
      console.log('initing');
      const self = this;
      const res = await getLanguageList();
      this.langList = res.data.languages;
      const { id } = self.$route.params;
      try {
        const res2 = await getContest(id);
        const { data } = res2;
        self.contest = data.contest;
        self.seeable = data.seeable;
        self.reason = data.reason;
      } catch (err) {
        console.log(err);
        self.$router.replace({ name: '404Page' });
      }
    },
    jumpToStatus() {
      const self = this;
      const routerResolve = self.$router.resolve({
        name: 'contestStatus',
        params: {
          id: self.contest.id,
        },
      });
      window.open(routerResolve.href, '_blank');
    },
    jumpToRank() {
      const self = this;
      const routerResolve = self.$router.resolve({
        name: 'contestRank',
        params: {
          id: self.contest.id,
        },
      });
      window.open(routerResolve.href, '_blank');
    },
    jumpToTeamRank() {
      const self = this;
      const routerResolve = self.$router.resolve({
        name: 'contestTeamRank',
        params: {
          id: self.contest.id,
        },
      });
      window.open(routerResolve.href, '_blank');
    },
  },

  beforeRouteUpdate(to, from, next) {
    console.log('beforeRouteUpdate!!');
    this.init();
    next();
  },
};
</script>
<style lang="scss" scoped>
</style>
