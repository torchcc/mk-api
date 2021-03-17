#!/usr/bin/env groovy
// Declarative //

def getLatestVersion(branch) {
    if (branch == 'release') {
        return 'RELEASE-LATEST'
    } else {
        return 'SNAPSHOT'
    }
}

def build(branch) {
    echo 'going to build branch ' + branch
    sh "go mod download"
    sh "go build -o app ."
    if (branch == 'test') {
        echo 'building test env docker image...'
        sh "docker build . -t t-mk-img -f ./deployment/test/Dockerfile"
        echo 'running test env docker container...'
        sh "docker rm -f t-mk-con && docker run -p 8081:8081 -d --name t-mk-con t-mk-img"
    }  else if (branch == 'release') {
        echo 'building prod env docker image...'
        sh "docker build . -t mk-img-prod -f ./deployment/prod/Dockerfile"
        echo 'running prod env docker container...'
        sh "docker rm -f p-mk-con && docker run -p 8071:8071 -d --name p-mk-con mk-img-prod"
    }
}

pipeline {
    agent any

    tools {
        go 'go-1.16'
    }
    environment {
        GO111MODULE = 'on'
        CGO_ENABLED = 0
        GOOS = 'linux'
        GOARCH = 'amd64'
        GOPROXY = 'https://goproxy.cn,direct'
        SERVICE_NAME = 'mk-api'
        TZ = 'Asia/Shanghai'
        scmVars = null
    }

    triggers {
        githubPush()
    }

    stages {
        stage('Prepare Env') {
            steps {
                echo 'Preparing Env...'
                // need to install workspace plugin
//                 cleanWs()
                checkout([$class: 'GitSCM', branches: [[name: '*/release']], doGenerateSubmoduleConfigurations: false, extensions: [], submoduleCfg: [], userRemoteConfigs: [[credentialsId: '47238156-6f3a-4339-9495-12d51b6c9577', url: 'git@github.com:Torchcc/mk-api.git']]])
                echo "checkout to path ${env.WORKSPACE}"
            }
        }
        stage('Build') {
            steps {
                echo "Running ${env.BUILD_ID} on ${env.JENKINS_URL}"
                build('release')
            }
        }
    }
    post {
        always {
            emailext(
                subject: '构建通知：${PROJECT_NAME} - Build # ${BUILD_NUMBER} -${BUILD_STATUS}!',
                body: '${FILE,path="email.html"}',
                to: 'troymm@163.com'
            )
        }
    }
}
