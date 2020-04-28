pipeline {
  agent any

  triggers {
    githubPush()
  }

  environment {
    VERSION="${env.BUILD_NUMBER}"
    GOOS="linux"
    GOARCH="amd64"
  }

  stages {
    stage('master: build and test') {
      when {
        branch 'master'
      }
      environment {
        ENV="master"
      }
      steps {
        slackSend (
          color: "#00F",
          message: "Starting build and test for <${env.BUILD_URL}|${env.JOB_NAME} #${env.BUILD_NUMBER}>"
        )

        sh '''
          make
        '''
      }
    }
  }

  post {
    success {
      slackSend (
        color: "#0F0",
        message: "<${env.BUILD_URL}|${env.JOB_NAME} ${env.BUILD_NUMBER}> success!"
      )
    }
    failure {
      slackSend (
        color: "#F00",
        message: "<${env.BUILD_URL}|${env.JOB_NAME} ${env.BUILD_NUMBER}> failed!"
      )
    }
  }
}

// vim: set expandtab ts=2 sw=2:
