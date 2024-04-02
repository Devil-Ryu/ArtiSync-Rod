# README

## About

This is the official Wails Vue template.

You can configure the project by editing `wails.json`. More information about the project settings can be found
here: https://wails.io/docs/reference/project-config

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.


## 待办清单
[] 系统配置：无头浏览器配置，超时时长(全局超时错误处理)，图片文件夹位置，是否打开软件时自动检测平台链接状况
[] CSDN: 增加错误处理，如长时间没有等待到响应的
[] Rodcontroller: 后续考虑多个平台使用一个rod，节省开支，增加速度