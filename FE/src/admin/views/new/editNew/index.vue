<template>
  <new :isEdit="true"></new>
</template>

<script>
import New from '../new/index.vue';

export default {
  name: 'adminEditNew',
  components: {
    New,
  },
  mounted() {
    console.log('mounted!!!!!!!!!!!!!!!');
  },
  activated() {
    console.log('activated!!!!!!!!!!!!!!!');
  },
};
</script>
