pipeline {
  agent any

  environment {
    IMAGE_NAME = "donation-backend"
    CONTAINER_NAME = "donation-backend"
  }

  stages {
    stage('Clone') {
      steps {
        git credentialsId: 'github-creds',
            url: 'git@github.com:Antonshepitko/donation-app.git'
      }
    }

    stage('Build Docker Image') {
      steps {
        sh 'docker build -t $IMAGE_NAME .'
      }
    }

    stage('Stop & Remove Old Container') {
        steps {
            sh '''
            if [ "$(docker ps -a -q -f name=donation-backend)" ]; then
                docker stop donation-backend
                docker rm donation-backend
            fi
            '''
        }
    }

    stage('Run New Container') {
      steps {
        sh '''
        if ! docker network ls --format '{{.Name}}' | grep -w donation-net > /dev/null; then
            docker network create donation-net
        fi
        docker run -d --name $CONTAINER_NAME \
            --network donation-net \
            -p 5000:5000 \
            $IMAGE_NAME
        '''
      }
    }
  }
}
