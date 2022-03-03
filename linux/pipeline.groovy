pipeline {
    agent any
    parameters {
        gitParameter name: 'branch',
        type: 'PT_BRANCH',
        branchFilter: 'origin/(.*)',
        defaultValue: 'master',
        selectedValue: 'DEFAULT',
        sortMode: 'ASCENDING_SMART',
        description: "选择需要的分支"
    }

    stages {
        stage('服务信息') {
            steps {
                sh 'echo 分支: $branch'
                sh 'echo 构建服务类型: ${JOB_NAME}-$type'
            }
        }

        stage('拉取代码') {
            steps {
                checkout([$class: 'GitSCM',
                          branches: [[name: '$branch']],
                          doGenerateSubmoduleConfigurations: false,
                          extensions: [],
                          submoduleCfg: [],
                          userRemoteConfigs: [[credentialsId: 'gitlab-cert', url: 'ssh://']]])
            }
        }

        stage('获取commit_id'){
            steps {
                echo '获取commit_id'
                git credentialsId: 'gitlab-cert', url: 'ssh://'
                script {
                    env.commit_id = sh(returnStdout: true, script: 'git rev-parse --short HEAD').trim()
                }
            }
        }

        stage('拉取配置文件') {
            steps {
                checkout([$class: 'GitSCM',
                          branches: [[name: '$branch']],
                          doGenerateSubmoduleConfigurations: false,
                          extensions: [[$class: 'RelativeTargetDirectory', relativeTargetDir: 'conf']],
                          submoduleCfg: [],
                          userRemoteConfigs: [[credentialsId: 'gitlab-cert', url: 'ssh://']]])
            }
        }

        stage('goctl版本检测') {
            steps {
                sh '/usr/local/bin/goctl -v'
            }
        }

        stage('Dockerfile Build') {
            steps {
                sh 'yes | cp -rf conf/${JOB_NAME}/${type}/${JOB_NAME}.yaml app/${JOB_NAME}/cmd/${type}/etc'
                sh 'cd app/${JOB_NAME}/cmd/${type} && /usr/local/bin/goctl docker -go ${JOB_NAME}.go && ls -l'
                script {
                    env.image = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}:${commit_id}').trim()
                }
                sh 'echo 镜像名称: ${image} && cp app/${JOB_NAME}/cmd/${type}/Dockerfile ./ && ls && docker build -t &{image} .'
            }
        }

        stage('上传到镜像仓库') {
            steps {
                sh 'docker login --username=${docker_name} --password=${docker-pwd} https://${docker_repo}'
                sh 'docker tag ${image} ${docker_repo}/go-zero-looklook/$image'
                sh 'docker push ${docker_repo}/go-zero-looklook/${image}'
            }
        }

        stage('部署到k8s') {
            steps {
                script {
                    env.deployYaml = sh(returnStdout: true, script: 'echo ${JOB_NAME}-${type}-deploy.yml').trim()
                    env.port = sh(returnStdout: true, script:'/root/port.sh ${JOB_NAME}-${type}').trim()
                }
                sh 'echo ${port}'

                sh 'rm -f ${deployYaml}'
                sh '/usr/local/bin/goctl kube deploy -secret docker-login -replicas 2 -nodePort 3${port} -requestCpu 200 -requestMem 50 -limitCpu 300 -limitMem 100 -name ${JOB_NAME}-${type} -namespace go-zero-looklook -image ${docker_repo}/${image} -o ${deployYaml} -port ${port} --home /root/template'
                sh '/usr/local/bin/kubectl apply -f ${deployYaml}'
            }
        }

        stage('Clean') {
            steps {
                sh 'docker rmi -f {image}'
                sh 'docker rmi -f ${docker_repo}/${image}'
                cleanWs notFailBuild: true
            }
        }
    }
}