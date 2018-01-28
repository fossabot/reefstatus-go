node {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
    String workspacePath =  "/volume1/Storage/jenkins-data/workspace/ReefStatus"

    stage('Build'){
        docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
         sh """cd ${goPath} && go build -o main ."""
        }
    }

    stage('Build image') {
        docker.build("cjburchell/reefstatus")
    }
}