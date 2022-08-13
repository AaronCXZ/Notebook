### Jenkinsfile

```jenkinsfile
pipeline {
	agent	{
		docker {
			image 'node:6-alpine'
			args	'-p 3000:3000 -p 5000:5000'
		}
	}
	environment {
		CI = 'true'
	}
	stages {
		stage('Build') {
			steps {
				sh 'npm install'
			}
		}
		stage('Test') {
			steps {
				sh './jenkins/scripts/test.sh'
			}
		}
		stage('Deliver for development') {
			when {
				branch 'dev'
			}
			steps {
				sh './jenkins/scripts/devliver-for-dev.sh'
				input message: 'Finished using the web site? (Click "Proceed" to continue)'
				sh './jenkins/scripts/kill.sh'
			}
		}
		stage('deploy for prod') {
			when {
				branch 'prod'
			}
			steps {
				sh ''
			}
		}
	}
}
```

Jenkins流水线是一套插件，它支持实现和集成continuous delivery（CD） pipeline到Jenkins

Jenkins流水线的定义被写在一个文本文件（jenkinsfile）中，该文件可以被提交到项目的源代码的控制仓库，这是’流水线即代码‘的基础。将CD流水线作为应用程序的一部分，像其他代码一样进行版本化和审查。可以获得以下好处：

- 自动地为所有分支创建流水线构建过程并拉取请求。
- 在流水线上代码复查/迭代（以及剩余的源代码）。
- 对流水线进行审计跟踪。
- 该流水线的真正的源代码可以被项目的多个成员查看赫然编辑。

Jenkinsfile能用两种语法进行编写：声明式和脚本化。通常都是两种写法相结合的流水线。

本质上，Jenkins是一个自动化引擎，它支持许多自动模式。流水线向Jenkins中添加了一组请打的工具，支持用例简单的持续集成到全面的CD流水线。通过对一系列的相关人物进行建模，用户可以利用流水线的很多特性：

- Code：流水线是在代码中实现的，通常会检查到源代码控制，使团队有编辑/审查和迭代他们的交付流水线的能力/
- Durable：流水线可以从Jenkins的主分支的计划内和计划外的重启中存活下来。
- Pausatile：流水线支持复杂的现实世界的CD请求，包括fork/join，循环，冰箱执行工作的能力。
- Extensible：流水线插件支持扩展到它的DSL的惯例和与其他插件集成的多个选项。

### 声明式流水线基础

pipeline块定义了整个流水线中完成的所有工作

 ```jenkinsfile
 pipeline {
 	agent any // 指定执行流水线的节点，any表示任意节点执行
 	stages {
 		stage('Build') { // 定义“Build”阶段
 			steps {
 				// 执行与“Build”阶段相关的步骤
 			}
 		}
 		stage('Test') { // 定义“Test”阶段
 			steps {
 				// 执行与“Test”阶段相关的步骤
 			}
 		}
 		stage('Deploy') { // 定义“Deploy”阶段
 			steps {
 				// 执行与“Deploy”阶段相关的步骤
 			}
 		}
 	}
 }
 ```

### 脚本化流水线基础

一个或多个node块在整个流水线中执行核心工作。限制了流水线在node块内的工作：

- 通过在Jenkins队列中添加一项来调度块中包含的步骤。节点上的执行器一空闲，该步骤就会执行。
- 创建一个工作区（特定为特定流水间建立的目录），其中工作可以在从源代码控制检出的文件上完成。

 ```jenkinsfile
 node { // 执行流水线的节点
 	stage('Build') { // 定义“Build”阶段，在脚本化流水线中实现stage块，可以清楚的显示Jenkins UI中的每个stage的任务子集
 		// 执行与“Build”阶段相关的步骤
 	}
 	stage('Test') { // 定义“Test”阶段
 		// 执行与“Test”阶段相关的步骤
 	}
 	stage('Deploy') { // 定义“Deploy”阶段
 		// 执行与“Deploy”阶段相关的步骤
 	}
 }
 ```

### 流水线语法

#### 节段

声明式流水线中的节段通常包含一个或多个执行或步骤。

##### 代理

agent部分指定了整个流水线或特定的部分，将会在Jenkins环境中执行的位置，该部分必须在pipeline的顶部被定义，但stage级别的使用是可选的。

##### 参数

为了支持作者可能有各种各样的用例流水线，agent部分支持一些不同类型的参数。这些参数在pipline的顶层或stage指令内部

###### any

在任何可用的代理上执行流水线或阶段

###### none

当pipeline块的顶部没有全局代理，该参数将会被分配到整个流水线的运行中并且每个stage部分都需要包含他自己的agent部分。

###### label

子啊提供了标签的Jenkins环境中可用的代理商执行流水线或阶段

###### node

agent { node {label 'labelName' } } 和 agent  { label 'labelName' }一样，但是node允许额外的选项

###### docker

使用给定的容器执行流水线或阶段。该容器将在预置的node上，或者匹配可选定义的label参数上，动态的供应来接受基于Docker的流水线。

docker也可以选择的接受args参数，该参数可能包含直接传递到docker run调用的参数以及alwaysPull选项，该选项强制docker pull，即使镜像名称已经存在。比如

```
agent {
	docker {
		image 'maven:3-alpine'
		label 'my-defined-label'
		args '-v /tmp:/tmp'
	}
}
```

###### dockerfile

执行流水线或阶段，使用从源代码库包含的Dockerfile构建的容器，jenkinsfile必须从多个分支流水线中加载，或者加载“Pipeline from SCM”。通常这是源代码仓库的根目录下的Dockerfile：agent { dockerfile true }。如果在另一个目录下构建Dockerfile，使用dir选项：agent { dockerfile  { dir 'someSubdir' } }。如果dockerfile有另一个名称可以使用filename选项指定文件名。可预约传递额外的参数到docker build ... 使用additionalBuildArgs选项提交，比如agent { dockerfile { additionalBuildArgs '--buiuld-arg foo=bar' } }。实例：

```
agent {
	dockerfile {
		filename 'Dockerfile.build'
		dir 'build'
		label 'my-defined-label'
		additionalBuildArgs '--build-arg version=1.0.2'
	}
}
```

##### 常见选项

有一些应用与两个或更多agent的实现的选项

###### label

一个字符串，该标签用于运行流水线或个别的stage，该选项对node/docker和dockerfile可用，node要求必须选择该选项

###### customWorkspace

一个字符串，在自定义工作区运行应用了agent的流水线或个别stage，而不是默认值。它既可以是一个相对路径，在这种情况下，自定义工作区会存在于节点工作区根目录下，或者一个绝对路径。比如：

```
agent {
	node {
		label 'my-defined-label'
		customWorkspace '/some/other/path'
	}
}
```

该选项对node，docker和dockerfile有用。

###### reuseNode

一个布尔值，默认true。如果是true，则流水线的顶层指定的节点上运行容器，在同样的工作区，而不是一个全新的节点上。对docker和dockerfile有用，并且只有当个别的stage的agent上才会有效。

##### 示例

```
pipeline {
	agent { docker 'maven:3-alpine' } // 在一个给定名称和标签的新建容器上执行定义的流水线中的所有步骤
	stages {
		stage('Ecample Build') {
			steps {
				sh 'mvn -B clean verify'
			}
		}
	}
}
```

```
pipeline {
	agent none // 在流水线顶层定义agent none确保an Exector没有被分配。使用agent node也会强制stage部分包含它自己的agent部分
	stages {
		stage('Ecample Build') {
			agent { docker 'maven::3-alpine'} // 使用镜像在一个新建的容器中执行该阶段的该步骤
			steps {
				echo 'Hello, Maven'
				sh 'mvn --version'
			}
		}
		stage('Ecample Test') {
			agent { docker 'openjdk:8-jre' } // 使用一个与之前阶段不同的镜像在一个新的容器中执行该阶段的该步骤
			steps {
				echo 'Hello, JDK'
				sh 'java -version'
			}
		}
	}
}
```

#### post

post部分定义一个或多个steps，这些阶段根据流水线或阶段的完成情况而运行（取决于流水线中post部分的位置）。post支持以下post-condition块中的其中之一：always，changed，failure，success，unstable和aborted。这些条件允许在post部分的步骤的执行取决于流水线或阶段的完成状态。

##### Conditions

###### always

无论流水线或阶段的完成状态如何，都允许post部分运行该步骤。

###### changed

只有当前流水线或阶段的完成状态与它之前的运行不同是，才允许在post部分运行该步骤

###### failure

只有当前流水线或阶段的完成状态为‘failure’，才允许在post部分运行该步骤，通常web UI是红色

###### success

只有当前流水线或阶段的完成状态为‘success’，才允许在post部分运行该步骤，通常web UI是蓝色或绿色

###### unstable

只有当前流水线或阶段的完成状态为‘unstable’，才允许在post部分运行该步骤，通常由于测试失败，代码违规等造成。通常web UI是黄色

###### aborted

只有当前流水线或阶段的完成状态为‘aborted’，才允许运行该步骤，通常由于流水线被手动的aborted，通常web UI是灰色

##### 示例

```
pipeline {
	agent any
	stages {
		stage('Example') {
			steps {
				echo 'Hello World'
			}
		}
	}
	post { // 按惯例post部分应该放在流水线的底部
		always {  // post-condition块包含与steps部分相同的steps
			echo 'I will always say Hello angain'
		}
	}
}
```

#### stages

包含一系列一个或多个stage指令，stages部分是流水线描述的大部分‘work’的位置，stages至少包含一个stage指令用于连续交付过程的每个离散部分，比如构建，测试，部署

##### 示例

```
pipleline {
	agent any
	stages { // stages部分通常会遵循诸如agent，options等指令
		stage('Example') {
			steps {
				echo 'Hello world'
			}
		}
	}
}
```

#### steps

steps部分在给定的stage指令中执行的定义了一系列的一个或多个steps

##### 示例

```
pipleline {
	agent any
	stages {
		stage('Example') {
			steps { // steps部分必须包含一个或多个步骤
				echo 'Hello world'
			}
		}
	}
}
```

#### 指令

##### enviroment

enviroment指令指定一个或多个健值对序列，该序列将被定义为所有步骤的环境变量，或者是特定于阶段的步骤，这取决于enviroment指令在流水线的位置。

该指令支持一个特殊的助手方法credentials()，该方法可用于在Jenkins环境中通过标识符访问预定义的凭证。

- 对于类型为‘Secret Text’的凭证，credentials()将确保指定的环境变量包含秘密文本的内容。
- 对于类型为‘SStandard username and password’的凭证，实现的环境变量指定为username:password，并且两个额外的环境变量被自动定义，分别为MYVARNAME_USR和MYUVARNAME_PSW

###### 示例

```
pipleline {
	agent any
	enviroment { // 顶层流水线使用enviroment指令将适用于流水线的所有步骤
		CC = 'clang'
	}
	stages {
		stage('Example') {
			enviroment {  // 在一个stage中定义的enviroment指令只会将给定的环境变量应用于stage中的步骤
				AN_ACCESS_KEY = credentials('my-prefined-secret-text') // 使用助手方法credentials()通过标识符获取Jenkins环境中预定义的凭证
			}
			steps {
				echo 'Hello world'
			}
		}
	}
}
```

#### options

options指令允许从流水线内部配置特定于流水线的选项。流水线提供了许多这样的选项，比如buildDiscarder，也可以由插件提供，比如timestamps。

##### 可用选项

###### buildDisacrder

为最近的流水线运行的特定数量保存组件和控制台输出。例如`options { buildDisacrder(logRotator(numTokeepStr: '1')) }`

###### disableConcurrentBuilds

不允许同时执行流水线，可被用来防止同时访问共享资源。例如` options { disableConcurrentBuilds() }`

###### overrideIndexTriggers

允许覆盖分支索引触发器的默认处理。如果分支索引在触发器在多个分支或组织标签中禁用，` options { overrideIndexTriggers(true) } ` 将只允许他们用于促工作，否则` options { overrideIndexTriggers(false) }`只会禁用该作业的分支索引触发器

###### skipDefaultChechout

在agent指令中，跳过从源代码控制中检出代码的默认情况

###### skipStagesAfterUnstable

一旦构建状态变得unstable，跳过该阶段

###### checkoutToSubdurectory

在工作空间的子目录中自动地执行源代码检出

###### timeout

设置流水线运行的超时时间，在此之后，Jenkins将终止流水线

###### retry

在失败时，重新尝试整个流水线的指定次数

###### timestamps

预谋所有流水线生成的控制台输出，与该流水线发出的时间一致

##### 示例

```
pipleline {
	agent any
	options {
		timeout(time: 1, unit: 'HOURS') // 指定一小时全局超时时间
	}
	stages {
		stage('Example') {
			steps {
				echo 'Hello world'
			}
		}
	}
}
```

##### 阶段选项

stage的options指令类似于流水线根目录上的options指令，但是stage级别的options只能包含retry，timeout，timestamps等步骤，或与stage相关的声明式选项

###### skipDefaultCheckout

在agent指令中跳过默认的从源代码检出代码的操作

###### timeout

此阶段的超时时间

###### retry

在失败时，重试此阶段的次数

###### timestamps

此阶段生成的所有控制它输出以及运行发出的时间一致

###### 示例

```
pipleline {
	agent any
	stages {
		stage('Example') {
			options {
				timeout(time: 1, unit: "HOURS") // 该阶段的超时时间是一小时
			}
			steps {
				echo 'Hello world'
			}
		}
	}
}
```

#### 参数

parameters指令提供了一个用户在触发流水线时应该提供的参数列表，这些用户指定的参数可以用过params对象提供给流水线步骤。

##### 可用参数

###### string

字符串类型的参数

###### booleanParam

布尔参数

##### 示例

```
pipleline {
	agent any
	parameters {
		string(name: 'PERSON', defaultValue: 'Mr Jenkins', description: 'Who should I say hello to?') // 定义一个name为PERSON，默认值为Mr Jenkins的参数
	}
	stages {
		stage('Example') {
			steps {
				echo "Hello ${params.PERSON}" // 获取用户输入参数
			}
		}
	}
}
```

#### 触发器

triggers指令定义了流水线被重新触发的自动化方法。对于集成了GitHub或BitBucket的流水线可能不需要triggers，因为给予web的集成很可能已经存在了。当前的触发器是cron，pollSCM和upstream

###### cron

接收cron央视的字符串来定义要重新触发流水线的常规间隔，比如` triggers { cron('H */4 * * 1-5') }`

###### pollSCM

接收cron样式的字符串来定义一个固定的间隔，在这个间隔中Jenkins会检查新的源代码更新，如果存在更改，流水线就会被重新触发

###### upstream

接收逗号分隔的字符串和阈值。当字符串中的任何作业以最小阈值结束时，流水线被重新触发，例如` triggers { upstream(upstreamPorjects: 'job1,job2', threshold: hudson.model.Result.SUCCESS) }`

##### 示例

```
pipleline {
	agent any
	triggers {
		cron('H */4 * * 1-5')
	}
	stages {
		stage('Example') {
			steps {
				echo 'Hello world'
			}
		}
	}
}
```

#### stage

stage指令在stages部分进行，每个stages应该包含一个。实际上流水线所有的工作都封装在一个或多个stage指令中

#### 工具

定义自动安装或防治PATH的工具的一部分，如果agent node指定，则忽略该操作

##### 支持工具

###### maven

###### jdk

###### gradle

#### input

stage的input指令允许你使用input step提示输入。在应用了options后，进入stage的agent或评估when条件前，stage将暂停，如果input被批准，stage将会继续。作为input提交的一部分的任何参数都将在环境变量中用于其他的stage

##### 配置项

###### message

必须的，这将在用户提交input时呈现给用户

###### id

input的可选标识，默认为stage的名称

###### ok

input表单上的ok按钮的可选文本

###### submitter

可选的以逗号分隔的用户列表或允许提交input的外部组名。默认允许任何用户

###### submitterParameter

环境变量的可断名称，如果存在，用submitter名称设置

###### parameters

提示提交者提供的一个可选的参数列表

##### 示例

```
pipleline {
	agent any
	stages {
		stage('Example') {
			input {
				message "Should we continue?"
				ok "Yes,we shoule."
				submitter "alice,bob"
				parameters {
					string(name: "PERSON", defaultValue: 'Mr Jenkins', description: 'Who should I say hello to?')
				}
			}
			steps {
				echo "Hello, ${PERSON}, nice to meet you."
			}
		}
	}
}
```

#### when



























