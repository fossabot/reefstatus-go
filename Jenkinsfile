node {
     stage('Clone repository') {
         /* Let's make sure we have the repository cloned to our workspace */
         checkout scm
     }

     String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
     String workspacePath =  "/data/jenkins/workspace/ReefStatus"

     stage('Build') {
       docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
           sh """cd ${goPath} && go build -o main ."""
          }
     }

     def app
     stage('Build image') {
          app = docker.build("cjburchell/reefstatus")
     }
 }