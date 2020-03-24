<template lang="pug">
  .content
    .content__main
      .one-main.has__pagination
        el-table(size="small",:data="tableData",v-loading="loading")
          el-table-column(label="排名", width="60")
            template(slot-scope="scope")
              span {{ (currentPage-1) * 50 + scope.$index + 1}}
          el-table-column(label="用户",width="70")
            template(slot-scope="scope")
                router-link(:to="{name:'userinfo',params:{id:scope.row.id}}")
                  .user__avatar__wrapper
                    img(:src="imgUrl(scope.row.avatar)",class="user__avatar")
          el-table-column
            template(slot-scope="scope")
              router-link(:to="{name:'userinfo',params:{id:scope.row.id}}")
                span {{`${scope.row.nick}`}}
          el-table-column(label="通过率", width="80")
            template(slot-scope="scope") {{calcRate(scope.row)}}
          el-table-column(label="解决",width="70",prop="solved")
          el-table-column(label="提交",width="70",prop="submit")

        el-pagination.user__pagination(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:pager-count="5",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total")
</template>

<script>
import { getRankList } from 'user/api/nologin';
import { setTimeout } from 'timers';
import { mapState } from 'vuex';

export default {
  data() {
    return {
      loading: false,
      currentPage: 1,
      perpage: 50,
      problemId: 0,
      tableData: [],
      total: 0,
    };
  },
  computed: {
    ...mapState({
      device: (state) => state.app.device,
    }),
  },
  mounted() {
    this.fetchData();
  },
  methods: {
    test(row) {
      console.log(row);
    },
    async fetchData() {
      const self = this;
      self.loading = true;
      try {
        const res = await getRankList(self.currentPage, self.perpage);
        console.log(res);
        setTimeout(() => {
          const { data } = res;
          self.tableData = data.data;
          self.total = data.total;
          self.loading = false;
        }, 200);
      } catch (err) {
        console.log(err);
      }
    },
    calcRate(row) {
      const rate = row.submit === 0 ? 0 : row.solved / row.submit;
      return `${Number(rate * 100).toFixed(2)}%`;
    },
  },
};
</script>

<style lang="scss" scoped>
.user__avatar__wrapper {
  img {
    width: 40px;
    height: 40px;
    border-radius: 40px;
  }
}
</style>
