pipeline {
    agent {
        docker {
            image 'golang' 
            args '-p 3000:3000' 
        }
    }
    stages {
        stage('Build') { 
            steps {
                sh 'go build -o service' 
            }
        }
    }
}