const layout = () => import('../views/Layout.vue');

export default [
  {
    path: '/',
    component: layout,
    children: [
      {
        path: '',
        name: 'index',
        meta: {
          keepAlive: true,
          title: '首页 - AHPUOJ',
        },
        component: () => import('../views/home'),
      },
      {
        path: 'findpass',
        name: 'findpass',
        meta: {
          keepAlive: true,
          title: '找回密码 - AHPUOJ',
        },
        component: () => import('../views/findpass'),
      },
      {
        path: 'resetpass',
        name: 'resetpass',
        meta: {
          keepAlive: true,
          title: '重设密码 - AHPUOJ',
        },
        component: () => import('../views/resetpass'),
      },
      {
        path: 'problemset',
        name: 'problemSet',
        meta: {
          keepAlive: true,
          title: '问题集 - AHPUOJ',
        },
        component: () => import('../views/problem_set'),
      },
      {
        path: 'issues',
        name: 'issueList',
        meta: {
          keepAlive: false,
          title: '讨论区 - AHPUOJ',
        },
        component: () => import('../views/issue_list'),
      },
      {
        path: 'issue/:id',
        name: 'issue',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/issue'),
      },
      {
        path: 'problem/:id/issues',
        name: 'problemIssueList',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/issue_list'),
      },
      {
        path: 'status',
        name: 'status',
        meta: {
          keepAlive: true,
          title: '评测机 - AHPUOJ',
        },
        component: () => import('../views/status'),
      },
      {
        path: 'contest/:id/status',
        name: 'contestStatus',
        meta: {
          keepAlive: true,
          title: '评测机 - AHPUOJ',
        },
        component: () => import('../views/status'),
      },
      {
        path: 'contest/:id/rank',
        name: 'contestRank',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/contest_rank'),
      },
      {
        path: 'contest/:id/teamrank',
        name: 'contestTeamRank',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/contest_team_rank'),
      },
      {
        path: 'problem/:id',
        name: 'problem',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/problem'),
      },
      {
        path: 'contest/:id/problem/:num',
        name: 'contestProblem',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/problem'),
      },
      {
        path: 'contests',
        name: 'contestList',
        meta: {
          keepAlive: true,
          title: '竞赛&作业 - AHPUOJ',
        },
        component: () => import('../views/contest_list'),
      },
      {
        path: 'series/:id',
        name: 'series',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/series'),
      },
      {
        path: 'serieses',
        name: 'seriesList',
        meta: {
          keepAlive: true,
          title: '系列赛 - AHPUOJ',
        },
        component: () => import('../views/series_list'),
      },
      {
        path: 'ranklist',
        name: 'ranklist',
        meta: {
          keepAlive: true,
          title: '排名 - AHPUOJ',
        },
        component: () => import('../views/ranklist'),
      },
      {
        path: 'contest/:id',
        name: 'contest',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/contest'),
      },
      {
        path: 'solution/:id',
        name: 'solution',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/solution'),
      },
      {
        path: 'account',
        name: 'account',
        meta: {
          keepAlive: false,
          title: '账号设置 - AHPUOJ',
        },
        component: () => import('../views/account_setting'),
      },
      {
        path: 'myreplys',
        name: 'myreplys',
        meta: {
          keepAlive: false,
          title: '查看回复 - AHPUOJ',
        },
        component: () => import('../views/myreplys'),
      },
      {
        path: 'userinfo/:id',
        name: 'userinfo',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/userinfo'),
      },
      // hack方法 只刷新路由
      {
        path: 'refresh',
        name: 'refresh',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/refresh'),
      },

      // 404路由
      {
        path: '404',
        name: '404Page',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/404'),
      },
      {
        path: '*',
        name: '404',
        meta: {
          keepAlive: false,
        },
        component: () => import('../views/404'),
      },
    ],
  },
];
