node('docker') {
    stage('Clone repository') {
        /* Let's make sure we have the repository cloned to our workspace */
        checkout scm
    }

    String goPath = "/go/src/github.com/cjburchell/reefstatus-go"

    stage('Build'){

        sh "ls -al"
        sh "pwd"
        sh 'echo WORKSPACE: $WORKSPACE'

        docker.image('golang:1.8.0-alpine').inside("-v ${pwd()}:${goPath}"){

         // Debugging
         sh 'echo GOPATH: $GOPATH'
         sh "ls -al ${goPath}"
         sh "cd ${goPath}"
         sh "pwd"

         sh """cd ${goPath} && go build -o main ."""
        }
    }

    stage('Build image') {
        /* This builds the actual image; synonymous to
         * docker build on the command line */

        docker.build("cjburchell/reefstatus")
    }
}