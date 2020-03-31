<template lang="pug">
  .content
    .content__main
      .siderbar
        ul.siderbar__item__list
          li
            el-button(size="mini",round,@click="handleSearchByResetConf(0)") 重置
          li
            .section__title 查找问题：
            .siderbar__searchbar__wrapper
              el-input(size="small",style="max-width:20em", placeholder="请输入问题名或ID", @keyup.enter.native="handleSearchByParam", v-model="queryParam", maxlength="20", clearable)
                el-button(slot="append" icon="el-icon-search", @click="handleSearchByParam")
            .tags__wrapper
          li
            .section__title 按难度检索：
            ul.button-list
                li
                  el-button(size="mini",round,:class="[level == -1?'is-active':'']",@click="handleSearchByLevel(-1)") 全部
                li
                  el-button(size="mini",round,:class="[level == 0?'is-active':'']",@click="handleSearchByLevel(0)") 简单
                li
                  el-button(size="mini",round,:class="[level == 1?'is-active':'']",@click="handleSearchByLevel(1)") 中等
                li
                  el-button(size="mini",round,:class="[level == 2?'is-active':'']",@click="handleSearchByLevel(2)") 困难
          li
            .section__title 按标签检索：
            ul.button-list
              li
                el-button(size="mini",round,:class="[tagId == -1?'is-active':'']",@click="handleSearchByTag(-1)") 全部
              template(v-for="tag in tags")
                li
                  el-button(size="mini",round,:class="[tagId == tag.id?'is-active':'']",@click="handleSearchByTag(tag.id)") {{tag.name}}
      .main.has__pagination
        el-pagination(style="float:left;",@current-change="fetchData",:current-page.sync="currentPage",:page-size="perpage",:pager-count="5",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total")
        el-table(size="small",:data="tableData",v-loading="loading",row-key="id")
          el-table-column(width="40")
            template(slot-scope="scope")
              svg-icon(name="ok",v-if="scope.row.status == 1")
              svg-icon(name="wrong",v-else-if="scope.row.status == 2")
          el-table-column(label="ID", prop="id", width="60")
          el-table-column(label="标题", min-width="160")
            template(slot-scope="scope")
              router-link(:to="{name:'problem',params:{id:scope.row.id}}") {{scope.row.title}}
          el-table-column(label="难度", min-width="60",align="center")
            template(slot-scope="scope")
              el-tag(size="small",v-if="scope.row.level==0",type="success",effect="dark") 简单
              el-tag(size="small",v-if="scope.row.level==1",type="warnning",effect="dark") 中等
              el-tag(size="small",v-if="scope.row.level==2",type="danger",effect="dark") 困难
          el-table-column(label="标签", min-width="160",align="center")
            template(slot-scope="scope")
              el-tag(size="small",v-for="tag in scope.row.tags",type="info",:key="tag.id",style="margin-left:3px;",effect="dark") {{tag.name}}
          el-table-column(label="通过率", min-width="80",align="center")
            template(slot-scope="scope") {{calcRate(scope.row)}}
          el-table-column(label="通过", prop="accepted", min-width="60",align="center")
          el-table-column(label="提交", prop="submit", min-width="60",align="center")
        el-pagination.user__pagination(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:pager-count="5",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total")
</template>

<script>
import { getProblemList, getAllTags } from "user/api/nologin";
import { mapState } from "vuex";

export default {
  name: "problemSet",
  data() {
    return {
      loading: false,
      currentPage: 1,
      perpage: 50,
      queryParam: "",
      tableData: [],
      total: 0,
      level: -1,
      tagId: -1,
      tags: []
    };
  },
  computed: {
    ...mapState({
      device: state => state.app.device
    })
  },
  async mounted() {
    this.tagId = this.$store.getters.tagId;
    console.log(this.tagId);
    this.fetchData();
    const res = await getAllTags();
    this.tags = res.data.tags;
  },
  activated() {
    if (this.$store.getters.tagId !== -1) {
      this.tagId = this.$store.getters.tagId;
      this.fetchData();
    }
    this.$store.dispatch("bus/resetTag");
  },
  methods: {
    async fetchData() {
      window.pageYOffset = 0;
      document.documentElement.scrollTop = 0;
      document.body.scrollTop = 0;
      this.loading = true;
      try {
        const res = await getProblemList(
          this.currentPage,
          this.perpage,
          this.queryParam,
          this.level,
          this.tagId
        );
        const { data } = res;
        this.tableData = data.data;
        this.total = data.total;
        this.loading = false;
      } catch (err) {
        console.log(err);
      }
      this.$store.dispatch("bus/resetTag");
    },
    handleSearchByResetConf() {
      this.level = -1;
      this.tagId = -1;
      this.queryParam = "";
      this.fetchData();
    },
    handleSearchByParam() {
      this.currentPage = 1;
      this.fetchData();
    },
    handleSearchByLevel(level) {
      this.currentPage = 1;
      this.level = level;
      this.fetchData();
    },
    handleSearchByTag(tagId) {
      this.currentPage = 1;
      this.tagId = tagId;
      this.fetchData();
    },
    calcRate(row) {
      const rate = row.submit === 0 ? 0 : row.accepted / row.submit;
      return `${(rate * 100).toFixed(2)}%`;
    }
  }
};
</script>
