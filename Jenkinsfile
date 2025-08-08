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
          docker stop $CONTAINER_NAME || true
          docker rm $CONTAINER_NAME || true
        '''
      }
    }

    stage('Run New Container') {
      steps {
        sh '''
          docker run -d --name $CONTAINER_NAME \
            --network donation-net \
            -p 5000:5000 \
            $IMAGE_NAME
        '''
      }
    }
  }
}
