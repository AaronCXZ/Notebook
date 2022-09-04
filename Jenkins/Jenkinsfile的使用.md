## Jenkinsfile的使用

#### 构建

```
Jekinsfile(Declarative Pipeline)
pipeline {
	agent any
	stages {
		stage('Build') {
			steps {
				sh 'make' # sh 步骤调用 make 命令，只有命令状态码返回零时才会继续。非零则会失败。
				archiveArtifacts artifacts: '**/target/*.jar', fingerprint: true # archiveArtifacts 捕获符合模式（'**/target/*.jar'）匹配的交付文件并将其保存到Jenkins master 节点以供后续获取
			}
		}
	}
}
```

#### 测试

```
Jekinsfile(Declarative Pipeline)
pipeline {
	agent any
	stages {
		stage('Test') {
			steps {
				sh 'make check || true' # 使用内连的 shell 条件（sh 'make check || true' 确保 sh 步骤的退出码总是零，使后面的 junit 步骤有机会捕获和处理测试报告
				junit '**/target/*.xml' #捕获并关联匹配的 JUnit XML 文件
			}
		}
	}
}

```

#### 部署

```
Jekinsfile(Declarative Pipeline)
pipeline {
	agent any
	stages {
		stage('Deploy') {
			when {
				expression {
					currenBuild.result == null || currentBuild.result == 'SUCCESS' # 流水线访问 currenBuild.result 变量确定是否有任何测试的失败
				}
			}
			steps {
				sh 'make publish'
			}
		}
	}
}
```

#### 使用环境变量

Jenkins流水线可以访问的环境变量列表记录在'${JENKINS_URL}/pipeline-syntax/globals#env'

- BUILD_ID：当前构建的ID，1.587+版本与BUILD_NUMBER相同
- BUILD_NUMBER：当前构建号
- BUILD_TAG：字符串'Jenkins-${JOB_NAME}-${BUILD_NUMBER}'。可以放到源码 / jar 等文件
- BUILD_URL：可以定位此次构建的结果的URL
- EXECTOR_NUMBER：用于识别执行当前构建的执行者的唯一编号
- JAVA_HOME：如果任务配置使用了特定的 JDK ，这个变量就是此 JDK 的 JAVA_HOME
- JENKINS_URL：Jenkins服务器的完成URL
- JOB_NAME：本次构建的雪梅名称
- NODE_NAME：运行本次构建的节点名称
- WORKSPACE：workspace的绝对路径

#### 设置环境环境

声明式使用environment指令，脚本式使用withEnv

#### 动态设置环境变量

```
Jenkinsfile(Declarative Pipeline)
pipeline {
	agent any # agent 必须设置在流水线的最高级。如果设置为 agent none 会失败
	environment {
		CC = """${sh (
				returnStdout: true,
				script: 'echo "clang"'
			)}""" # 使用 returnStdout 时，返回的字符串末尾会追加一个空格，可以使用 .trim() 将其移除
	// 使用 returnStatus
	EXIT_STATUS= """${sh(
			returnStatus: true,
			script: 'exit 1'
		)}"""
	}
	stages {
		stage('Example') {
			envirinment {
				DEBUG_FLAGS = '-g'
			}
			strps {
				sh 'printenv'
			}
		}
	}
}
```

#### 处理凭据













