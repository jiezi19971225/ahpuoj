const bus = {
  namespaced: true,
  state: {
    // 标签检索题目
    tagId: -1,

    // 条件检索评测记录
    solutionQueryParam: '',
    solutionUserNick: '',
    solutionLanguage: -1,
    solutionResult: -1,
  },
  mutations: {
    SET_TAG: (state, tagId) => {
      state.tagId = tagId;
    },
    RESET_TAG: (state) => {
      state.tagId = -1;
    },
    SET_SOLUTION_FILTER: (state, {
      queryParam, nick, language, result,
    }) => {
      state.solutionQueryParam = queryParam || '';
      state.solutionUserNick = nick || '';
      state.solutionLanguage = language || -1;
      state.solutionResult = result || -1;
    },
    RESET_SOLUTION_FILTER: (state) => {
      state.solutionQueryParam = '';
      state.solutionUserNick = '';
      state.solutionLanguage = -1;
      state.solutionResult = -1;
    },
  },
  actions: {
    setTag({ commit }, tagId) {
      commit('SET_TAG', tagId);
    },
    resetTag({ commit }) {
      commit('RESET_TAG');
    },
    setSolutionFilter(
      { commit },
      {
        queryParam = 0, nick = '', language = -1, result = -1,
      },
    ) {
      commit('SET_SOLUTION_FILTER', {
        queryParam,
        nick,
        language,
        result,
      });
    },
    resetSolutionFilter({ commit }) {
      commit('RESET_SOLUTION_FILTER');
    },
  },
};

export default bus;
