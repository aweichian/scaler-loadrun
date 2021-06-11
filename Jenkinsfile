pipeline {
    agent {node 'slave-node06'}

    environment {
        USER_EMAIL = "renwei.qian@baozun.com"  //发送邮件的地址（一般为项目owner，按实际修改）
        USER = "renwei.qian"  //发送邮件的邮箱前缀（一般为项目owner，按实际修改）
        SFTP_SECRET_ACCESS_KEY = credentials('jenkins-sftp-secret-access-key')  //SFTP的秘钥（发包使用，无需修改）
        SFTP_SERVER = "10.101.6.87"  //SFTP的地址（发包使用，无需修改）
        VALIDATE_URL = "http://bee-backend.baozun.com/ci_validate/"  //SFTP的地址（发包使用，无需修改）
        REMOTE_DIR = "/upload"  //SFTP的目录（发包使用，无需修改）
        LOCAL_DIR= "$WORKSPACE/ecs-ofa-service-impl/target"  //生成包的位置，其中$WORKSPACE不需要修改，ecs-ofa-service-impl/target按照实际的应用目录填写
        PACKAGE_NAME = "scaler-loadrun"  //需要发布的包名，按实际修改，dubbo包名写.tar.gz包，tomcat包名写war包，springboot包名写jar包
        PROJECT_NAME = "devops" //项目名，按实际修改
        APP_NAME = "scaler-loadrun"  //应用名，按实际修改，如有多个，在括号内添加，以空格间隔(例如"ofa-service-a ofa-service-b ofa-service-c")
        ENV_NAME = "sit"  //发布的环境，根据自己有几个环境进行删减，但是环境名就sit uat sandbox prod这四个
        HARBOR_ADDR = "ic-harbor.baozun.com"  //Harbor地址，无需修改
        DOCKER_NAME = "docker/Dockerfile"  //Dockerfile文件位置，按实际修改
		COVERAGE = "50"  //单元测试覆盖率，标准为50,新接入应用可适当降低
    }

    options {
        buildDiscarder(logRotator(numToKeepStr: '5', artifactNumToKeepStr: '5'))  //保留历史记录，无需修改
    }

    //pipeline运行结果通知给触发者，无需修改
    post {
        success {
            script {
                wrap([$class: 'BuildUser']) {
                    emailext body: '$DEFAULT_CONTENT', recipientProviders: [developers()], mimeType: 'text/html', subject: '$DEFAULT_SUBJECT', to: "$USER_EMAIL"
                }
            }
        }
        failure {
            script {
                wrap([$class: 'BuildUser']) {
                    emailext body: '$DEFAULT_CONTENT', recipientProviders: [developers()], mimeType: 'text/html', subject: '$DEFAULT_SUBJECT', to: "$USER_EMAIL"
                    }
            }

        }
        unstable {
            script {
                wrap([$class: 'BuildUser']) {
                    emailext body: '$DEFAULT_CONTENT', recipientProviders: [developers()], mimeType: 'text/html', subject: '$DEFAULT_SUBJECT', to: "$USER_EMAIL"
                }
            }
        }
        aborted {
            script {
                wrap([$class: 'BuildUser']) {
                    emailext body: '$DEFAULT_CONTENT', recipientProviders: [developers()], mimeType: 'text/html', subject: '$DEFAULT_SUBJECT', to: "$USER_EMAIL"
                }
            }
        }
    }

	stages {
		stage('静态检查') {
            when { anyOf{branch 'testing';branch 'ci'} }  //哪些分支需要单元测试就可直接修改或者添加（例如添加release分支，在branch 'ci'后添加;branch 'release'）
            steps {
                timeout(time: 20, unit: 'MINUTES') {  //设置超时时间，无需修改
                    withSonarQubeEnv('sonarserver') {
					    //收集jacoco.exec文件，如果在当前项目生成target目录，需修改为$WORKSPACE/target/jacoco.exec，即不需要**
                    	sh '''dirSrc=`ls -m $WORKSPACE/**/target/jacoco.exec|tr -d " "`
                    	mvn -f ${POM_NAME} sonar:sonar -Dsonar.login=cicd -Dsonar.password=cicd -Dsonar.core.codeCoveragePlugin=jacoco -Dsonar.jacoco.reportPaths="$dirSrc" -Dsonar.dynamicAnalysis=reuseReports'''  //sonar扫描，maven的sonar扫描，如果是gradle需要配置sonar-scanner，并且外面需添加sonar.properties文件
                    }

                    script {
                        timeout(5) {
                            //利用sonar scaler-loadrun功能通知pipeline代码检测结果，未通过质量阈，pipeline将会fail
                            def qg = waitForQualityGate()
                            echo "${qg.status}"
                            if (qg.status != 'OK') {
                                echo "${qg.status}"
                                error "未通过Sonarqube的代码质量阈检查，请及时修改！failure: ${qg.status}"
                            }
                        }
                    }
                }
            }
        }
        stage('打包&上传镜像') {
            when { anyOf{branch 'master';branch 'ci'} }  //哪些分支需要打包上传镜像就可直接修改或者添加（例如添加release分支，在branch 'ci'后添加;branch 'release'）
            steps {
                timeout(time: 20, unit: 'MINUTES') {
                    sh '''
                        export GO111MODULE=on
                        export GOPROXY=https://goproxy.io
                        go build main.go

                        appname=()
                        dockername=()
                        envname=()
                        for app in $APP_NAME;do
                            appname=(${appname[@]} $app)
                        done
                        for docker in $DOCKER_NAME;do
                            dockername=(${dockername[@]} $docker)
                        done
                        for e_name in $ENV_NAME;do
                            envname=(${envname[@]} $e_name)
                        done
                        appname_len=${#appname[@]}
						envname_len=${#envname[@]}
                        for ((i = 0; i < appname_len; i++));do
                            for ((j = 0; j < envname_len; j++));do
                                docker images  --filter="reference=${HARBOR_ADDR}/${envname[$j]}/${PROJECT_NAME}_${appname[$i]}:*" -q | xargs --no-run-if-empty docker rmi --force
                                docker build --no-cache -t ${HARBOR_ADDR}/${envname[$j]}/${PROJECT_NAME}_${appname[$i]}:${GIT_COMMIT:0:7} -f  ${dockername[$i]} ./
                                docker push ${HARBOR_ADDR}/${envname[$j]}/${PROJECT_NAME}_${appname[$i]}:${GIT_COMMIT:0:7}
                            done
                        done
                        docker tag ${HARBOR_ADDR}/${ENV_NAME}/${PROJECT_NAME}_${appname}:${GIT_COMMIT:0:7} ${HARBOR_ADDR}/prod/${PROJECT_NAME}_${appname}:${GIT_COMMIT:0:7}
                        docker push ${HARBOR_ADDR}/prod/${PROJECT_NAME}_${appname}:${GIT_COMMIT:0:7}
					'''
                    }
              }
        }
    }
}
