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
    sh "go mod tidy"
    sh "go build -o app ."
    if (branch == 'release') {
        echo 'deploying...'
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
                cleanWs()
                checkout scm
                echo "checkout to path ${env.WORKSPACE}"
            }
        }

        stage('Build') {
            steps {
            echo "Running ${env.BUILD_ID} on ${env.JENKINS_URL}"
            sh 'mvn clean package'
            }
        }
        stage('publish project') {
            steps {
                deploy adapters: [tomcat9(credentialsId: '8c2f7e52-3591-4460-84fc-64ac1c1482ad', path: '', url: 'http://106.53.124.190:7080')], contextPath: 'web_demo_pipeline', war: 'target/*.war'
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
