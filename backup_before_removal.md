# 图像识别模块移除前备份记录

## 备份时间
2026-03-13

## 当前Git状态
- 分支：feature/remove-image-recognition
- 基于：main分支
- 未提交更改：3个修改文件

## 图像识别模块文件清单

### 后端文件：
1. `server/common/image/image_recognizer.go` - 222行
2. `server/controller/image/image.go` - 40行  
3. `server/service/image/image.go` - 38行
4. `server/router/Image.go` - 10行

### 前端文件：
1. `client/src/views/ImageRecognition.vue` - 462行

### 相关配置：
1. `server/router/router.go` - 包含图像路由组注册
2. `client/src/router/index.js` - 包含前端路由配置
3. `client/src/views/Menu.vue` - 包含菜单入口

## 依赖库状态
- ONNX Runtime: github.com/yalue/onnxruntime_go v1.23.0
- 图像处理库: golang.org/x/image v0.32.0

## 路由状态
- 后端路由：已注释，但路由组注册仍存在
- 前端路由：正常启用

## 风险控制
- 已创建专用分支进行移除操作
- 可随时回滚到当前状态