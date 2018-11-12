pipeline {
  agent {
    node {
      label 'master'
    }
  }
  parameters {
    string(name: 'distName', defaultValue: 'nap-syslog' , description: '镜像名')
    string(name: 'distNumber', defaultValue: "${BRANCH_NAME}", description: '镜像版本')
    string(name: 'distRepository', defaultValue: 'http://mvn.sky-cloud.net:8081/repository/nap-statics/nap-syslog/', description: '静态文件发布位置')
  }
  stages {
    stage('Build') {
      steps {
        script {
           docker.withRegistry('http://hub.sky-cloud.net', 'docker-image-builder') {
            docker.image('hub.sky-cloud.net/cicd/gobuilder:v3.0').inside {
              sh 'cp -r ./.  $GOPATH/src/github.com/ekanite/ekanite; cd $GOPATH/src/github.com/ekanite/ekanite; go get ./...'
              sh 'cd cmd/ekanited/ ; go build .'
              def PACKAGE_NAME = "${params.distName}.${params.distNumber}_build-${BUILD_NUMBER}"
              sh "cd /go/src/github.com/ekanite/ ; tar -zcvf ${PACKAGE_NAME}.tar.gz ./ekanite ; curl --fail -u admin:r00tme --upload-file ${PACKAGE_NAME}.tar.gz ${params.distRepository}"
            }
          }
        }
      }
    }
    stage('build-dockerimage') {
      steps {
        script {
          def PACKAGE_HOME = "${params.distRepository}${params.distName}.${params.distNumber}_build-${BUILD_NUMBER}"
          sh "curl -u admin:r00tme -O ${PACKAGE_HOME}.tar.gz; tar -zxvf ${params.distName}.${params.distNumber}_build-${BUILD_NUMBER}.tar.gz ; cd ekanite"
          sh "docker login -u docker-image-builder -p sky-cloud@SZ2018 hub.sky-cloud.net ; docker build -t hub.sky-cloud.net/nap/${params.distName}:${params.distNumber}_build-${BUILD_NUMBER} . ; docker push hub.sky-cloud.net/nap/${params.distName}:${params.distNumber}_build-${BUILD_NUMBER}"
        }
      }
    }
  }
}

