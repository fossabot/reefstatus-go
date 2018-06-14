node {
     stage('Clone repository') {
         /* Let's make sure we have the repository cloned to our workspace */
         checkout scm
     }

     String dockerImage = "cjburchell/reefstatus"
     String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
     String workspacePath =  """${env.WORKSPACE}"""

     stage('Build Frontend') {
      docker.image('node:10-alpine').inside("-v ${workspacePath}:${goPath}"){
                sh """sudo npm install -g @angular/cli"""
                sh """cd ${goPath}/frontend/ReefStatus && npm install"""
                sh """cd ${goPath}/frontend/ReefStatus && ng build --prod"""
               }
     }

     stage('Build Server') {
       docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
           sh """cd ${goPath} && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main"""
          }
     }

     stage('Build image') {
          docker.build("${dockerImage}").tag('latest')
     }

     stage ('Docker push') {
               docker.withRegistry('https://390282485276.dkr.ecr.us-east-1.amazonaws.com', 'ecr:us-east-1:redpoint-ecr-credentials') {
                docker.image("${dockerImage}").push('latest')
              }
     }
 }