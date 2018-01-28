node {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
    String workspacePath =  "/volume1/Storage/jenkins-data/workspace/ReefStatus"

    stage('Test') {
        docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){

            sh """cd ${goPath}"""
            sh 'go list ./...'
            sh 'go list ./... | grep -v /vendor/ > projectPaths.txt'
            def paths = sh returnStdout: true, script: """awk '"_/var/jenkins_home/workspace/ReefStatus"\$0="./src/github.com/cjburchell/reefstatus-go"\$0' projectPaths.txt"""

            sh 'echo paths: ${paths}'

            sh """cd /go && go tool vet ${paths}"""

            sh """cd /go && golint ${paths}"""

            sh """cd /go && go test -race -cover ${paths}"""
        }
    }

    stage('Build'){
        docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
         sh """cd ${goPath} && go build -o main ."""
        }
    }

    stage('Build image') {
        docker.build("cjburchell/reefstatus")
    }
}