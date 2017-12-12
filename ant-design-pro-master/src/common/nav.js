import dynamic from 'dva/dynamic';

// wrapper of dynamic
const dynamicWrapper = (app, models, component) => dynamic({
  app,
  models: () => models.map(m => import (`../models/${m}.js`)),
  component
});

// nav data
export const getNavData = app => [
  {
    component: dynamicWrapper(app, [
      'user', 'login'
    ], () => import ('../layouts/Dashboard')),
    layout: 'BasicLayout',
    name: '首页',
    path: '/',
    children: [
      {
        name: '站点设置',
        icon: 'dashboard',
        path: 'admin',
        children: [
          {
            name: '当前状态',
            path: 'status',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '基本信息',
            path: 'info',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '作者',
            path: 'author',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '搜索引擎',
            path: 'seo',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '邮件发送',
            path: 'smtp',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '国际化',
            path: 'locales',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '链接管理',
            path: 'links',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '卡片管理',
            path: 'cards',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '友情链接',
            path: 'friend-links',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '用户列表',
            path: 'users',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }
        ]
      }, {
        name: '交流论坛',
        path: 'forum',
        icon: 'form',
        children: [
          {
            name: '文章列表',
            path: 'articles',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '标签列表',
            path: 'tags',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '评论列表',
            path: 'comments',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }
        ]
      }, {
        name: '问卷调查',
        path: 'survey',
        icon: 'table',
        children: [
          {
            name: '设计问卷',
            path: 'forms',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '统计结果',
            path: 'records',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }
        ]
      }, {
        name: '附件管理',
        path: 'attachments',
        icon: 'cloud-upload-o',
        children: [
          {
            name: '文件上传',
            path: 'upload',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '我的文档',
            path: 'list',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }
        ]
      }, {
        name: '个人信息',
        path: 'users',
        icon: 'user',
        children: [
          {
            name: '日志列表',
            path: 'logs',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '基本信息',
            path: 'profile',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }, {
            name: '修改密码',
            path: 'change-password',
            component: dynamicWrapper(app, [], () => import ('../routes/nut/Home'))
          }
        ]
      }
    ]
  }, {
    component: dynamicWrapper(app, [], () => import ('../layouts/UserLayout')),
    layout: 'UserLayout',
    children: []
  }, {
    component: dynamicWrapper(app, [], () => import ('../layouts/BlankLayout')),
    layout: 'BlankLayout',
    children: {
      name: '使用文档',
      path: 'http://localhost:8080',
      target: '_blank',
      icon: 'book'
    }
  }
];
