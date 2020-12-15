## Flutter项目开发规范

#### 目录结构：

`assets`目录里再进行分类：如images、audios、videos、fonts等等

`lib`文件夹：

`base`：基础通用组件

`config`： 常用的全局配置（常量等）

`utils`：工具类

`widgets`：轮子

业务模块代码, 以业务模块分别为文件夹，细分为`model`， `pages`,`vm`

 ....

##### 命名规范

`UpperCamelCase`: 一般用在类名、注解、枚举、typedef和参数的类型

`lowerCamelCase`: 一般用在类成员、变量、方法名、参数命名等命名上

`lowercase_with_underscores`:一般用在命名库（libraries）、包(packages)、目录（directories）和源文件(source files)上

#### 代码规范

1. 不要出现三元运算符嵌套
2. 判空兼容
3. 使用///来表示，并且注释会出现生成到文档里。一般我们可以使用文档注释来注释类成员和类型、方法、参数、类、变量常量











​				



