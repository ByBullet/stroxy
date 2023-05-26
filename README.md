# stroxy

## 构建
1. 安装python3
2. 执行 cd .../stroxy && python build.py
3. 编译结果放在 .../stoxy/build/

## 运行
双击或命令行运行`stroxy`

## 开发约定
1. 为了让path_util.go中的GetCurrentAbPath()能适配环境，项目中不能出现名称为stroxy的文件夹或文件
2. 所有的TODO用大写
3. 要调用某些工具时去util包看看，别重复造轮子
4. 为了避免循环依赖，各个包之间的协助在boot包中处理