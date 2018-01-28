node {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    String goPath = "/go/src/github.com/cjburchell/reefstatus-go"
    String workspacePath =  "/volume1/Storage/jenkins-data/workspace/ReefStatus"

    stage('Test') {
        docker.image('golang:1.8.0-alpine').inside("-v ${workspacePath}:${goPath}"){
            sh """cd /go/src && go list ./..."""
            sh """cd ${goPath} && go list ./..."""
            def projectPaths = sh 'cd /go/src && go list ./... | grep -v /vendor/'
            def paths = sh returnStdout: true, script: """awk '\$0="./src/"\$0' ${projectPaths}"""

            sh 'echo paths: ${paths}'

            sh """cd ${goPath} && go tool vet ${paths}"""

            sh """cd ${goPath} && golint ${paths}"""

            sh """cd ${goPath} && go test -race -cover ${paths}"""
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