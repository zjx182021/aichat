# ai-chat-web 

## node 安装
### windows
1. 上[官网](https://nodejs.org/en)下载18.16.0 LTS版本
2. 查看node 是否安装成功
```
node -v
npm -v
```
3. 安装pnpm
```
npm install pnpm  -g
```

### ubuntu
1. 设置 apt 源，设置后可查看/etc/apt/sources.list.d/nodesource.list 文件
```
curl -sL https://deb.nodesource.com/setup_18.x | sudo -E bash -
```
2. 安装nodejs
```
sudo apt-get install -y nodejs
```
3. 验证
```
node -v
npm -v
```
4. 安装pnpm
```
sudo npm install pnpm -g
```

## 编译运行
1. 依赖安装
``` 
pnpm bootstrap
```
2. 本地运行
```
pnpm dev
```
3. 打包发布版本
```
pnpm build-only
```

## 提交代码的规则
```
* commitlint 规则是指在提交代码时要遵循的规范，常见的 commitlint 规则如下：

* type：用于说明 commit 的类型，例如 feat（新功能）、fix（修复 bug）、docs（文档更新）、style（样式修改）、refactor（重构代码）等。

* scope：用于说明 commit 影响的范围，例如组件、模块、页面等。

* subject：用于简短地描述 commit 的内容，建议不超过 50 个字符。

* body：用于详细描述 commit 的改动内容，可以分成多行。

* footer：用于关闭 issue 或者添加相关链接等信息。

* 长度限制：commit message 不应该过长，一般不超过 72 个字符。
```