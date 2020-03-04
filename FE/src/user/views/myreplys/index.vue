<template lang="pug">
  .content
    .content__main
      .one-main
        h1.content__panel__title 回复列表
        div(v-if="replys&&replys.length==0")
          .reply__box
            p 还没有人回复你
        div(v-else)
          template(v-for="item,index in replys")
            .reply__box
              router-link(:to="{name:'userinfo',params:{id:item?item.user.id:0}}") {{item.user.username}}
              span 在帖子
              router-link(:to="{name:'issue',params:{id:item?item.issue_id:0}}") {{item.issue_title}}
              span 中回复了你
              br
              .reply__content(v-html="calcContent(item.content)")
        el-pagination.tal.mt20(@current-change="fetchData",:current-page.sync="currentPage",background,
        :page-size="perpage",:layout="'prev, pager, next'+(device=='desktop'?',jumper':'')",:total="total",:small="device === 'mobile'")
</template>

<script>
import TinymceEditor from 'common/components/tinymce_editor.vue';
import EventBus from 'common/eventbus';
import { getMyReplys } from 'user/api/user';
import { mapState } from 'vuex';

export default {
  components: {
    TinymceEditor,
  },
  data() {
    return {
      currentPage: 1,
      dialogFormVisible: false,
      perpage: 20,
      replys: [],
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
    async fetchData(resetScroll) {
      if (resetScroll !== false) {
        window.pageYOffset = 0;
        document.documentElement.scrollTop = 0;
        document.body.scrollTop = 0;
      }
      const self = this;
      try {
        const res = await getMyReplys(self.currentPage, self.perpage);
        const { data } = res;
        console.log(res);
        self.replys = data.replys;
        self.total = data.total;
      } catch (err) {
        console.log(err);
      }
    },
    calcContent(content) {
      return content.length <= 100 ? content : `${content.substr(0, 100)}...`;
    },
  },
};
</script>

<style lang="scss" scoped>
.reply__box {
  background: $--color-level15;
  position: relative;
  margin-top: 0.2rem;
  padding: 0 0.1rem 0.1rem 0.1rem;
  border-bottom: 1px solid $--color-level13;
  text-align: left;

  .reply__content {
    border-radius: 5px;
    background: $--color-level14;
    padding: 0.1rem;
    position: relative;
    box-sizing: border-box;
    min-height: 60px;
    text-align: left;
    font-size: 16px;
  }
}
</style>
