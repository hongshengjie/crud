const path = require('path');
var htmlWebpackPlugin = require('html-webpack-plugin');
module.exports = {
  entry: './src/index.js',
  devtool: 'inline-source-map',
  mode:'development',
  output: {
    filename: 'bundle.js',
    path: path.resolve(__dirname, 'dist')
  },
  plugins:[ // 添加plugins节点配置插件
    new htmlWebpackPlugin({
        //template:path.resolve(__dirname, './src/index.html'),//模板路径
        filename:'index.html'//自动生成的HTML文件的名称
    })
]
};
