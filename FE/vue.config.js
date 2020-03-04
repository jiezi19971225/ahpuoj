const path = require('path');

module.exports = {
  lintOnSave: false,
  productionSourceMap: false,
  chainWebpack: (config) => {
    config.module.rule('pug')
      .test(/\.pug$/)
      .use('pug-html-loader')
      .loader('pug-html-loader')
      .end();

    const svgRule = config.module.rule('svg');
    // 清除已有的所有 loader。
    // 如果你不这样做，接下来的 loader 会附加在该规则现有的 loader 之后。
    svgRule.uses.clear();
    svgRule
      .test(/\.svg$/)
      .include.add(path.resolve(__dirname, './src/common/assets/icons'))
      .end()
      .use('svg-sprite-loader')
      .loader('svg-sprite-loader')
      .options({
        symbolId: 'icon-[name]',
      });
    const fileRule = config.module.rule('file');
    fileRule.uses.clear();
    fileRule
      .test(/\.svg$/)
      .exclude.add(path.resolve(__dirname, './src/common/assets/icons'))
      .end()
      .use('file-loader')
      .loader('file-loader');
  },
  css: {
    loaderOptions: {
      scss: {
        prependData: '@import "~@/common/styles/variables.scss";',
      },
    },
  },
  pages: {
    user: {
      // page 的入口
      entry: 'src/user/main.js',
      // 模板来源
      template: 'public/index.html',
      // 在 dist/index.html 的输出
      filename: 'index.html',
      // 当使用 title 选项时，
      // template 中的 title 标签需要是 <title><%= htmlWebpackPlugin.options.title %></title>
      title: 'AHPUOJ',
      // 在这个页面中包含的块，默认情况下会包含
      // 提取出来的通用 chunk 和 vendor chunk。
      // chunks: ['chunk-vendors', 'chunk-common', 'index'],
    },
    admin: {
      // page 的入口
      entry: 'src/admin/main.js',
      // 模板来源
      template: 'public/index.html',
      // 在 dist/index.html 的输出
      filename: 'admin/index.html',
      // 当使用 title 选项时，
      // template 中的 title 标签需要是 <title><%= htmlWebpackPlugin.options.title %></title>
      title: 'AHPUOJ 后台管理系统',
      // 在这个页面中包含的块，默认情况下会包含
      // 提取出来的通用 chunk 和 vendor chunk。
      // chunks: ['chunk-vendors', 'chunk-common', 'index'],
    },
  },
  configureWebpack: (config) => {
    if (process.env.NODE_ENV === 'production') {
      config.devtool = 'source-map';
    }
    return {
      resolve: {
      // 引入路径是不用写对应的后缀名
        extensions: ['.js', '.vue', '.json'],
        alias: {
          '@': path.resolve(__dirname, 'src'),
          user: path.resolve(__dirname, 'src/user'),
          admin: path.resolve(__dirname, 'src/admin'),
          common: path.resolve(__dirname, 'src/common'),
        },
      },
    };
  },
  devServer: {
    open: true,
    hot: true,
    host: '127.0.0.1',
    port: '8888',
    disableHostCheck: true,
    historyApiFallback: {
      rewrites: [
        {
          from: /^\/$/,
          to: '/index.html',
        },
        {
          from: /^\/admin/,
          to: '/admin/index.html',
        },
      ],
    },
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        ws: false,
        changOrigin: true,
        secure: false,
      },
      '/upload': {
        target: 'http://localhost',
        ws: false,
        changOrigin: true,
        secure: false,
      },
    },
  },
};
